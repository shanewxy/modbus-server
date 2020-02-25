package cmd

import (
	"math/rand"
	"time"

	"github.com/goburrow/serial"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tbrandon/mbserver"
	ctrl "sigs.k8s.io/controller-runtime"
)

var device rtuDevice

func NewCommand() *cobra.Command {
	var c = &cobra.Command{
		Use:  "rtu",
		Long: "test device rtu",
		RunE: func(cmd *cobra.Command, args []string) error {
			device.Run()
			return nil
		},
	}

	c.Flags().StringVarP(&device.address, "address", "a", "/dev/ttyS0", "rtu listening ip address of the modbus server")
	c.Flags().StringVarP(&device.parity, "parity", "p", "E", "Parity: N - None, E - Even, O - Odd (default E)")
	c.Flags().IntVarP(&device.baudRate, "baudRate", "b", 19200, "rtu baudRate of the modbus server")
	c.Flags().IntVarP(&device.dataBits, "dataBits", "d", 8, "Data bits: 5, 6, 7 or 8 (default 8)")
	c.Flags().IntVarP(&device.stopBits, "stopBits", "s", 1, "Stop bits: 1 or 2 (default 1)")

	return c
}

type rtuDevice struct {
	// Device path (/dev/ttyS0)
	address string
	// Baud rate (default 19200)
	baudRate int
	// Data bits: 5, 6, 7 or 8 (default 8)
	dataBits int
	// Stop bits: 1 or 2 (default 1)
	stopBits int
	// Parity: N - None, E - Even, O - Odd (default E)
	// (The use of no parity requires 2 stop bits.)
	parity string
	// Read (Write) timeout.
}

func (r *rtuDevice) Run() {
	server := mbserver.NewServer()
	config := &serial.Config{
		Address:  r.address,
		BaudRate: r.baudRate,
		DataBits: r.dataBits,
		StopBits: r.stopBits,
		Parity:   r.parity,
	}
	if err := server.ListenRTU(config); err != nil {
		logrus.Error("Fail to start server", err)
		return
	}

	tick := time.NewTicker(time.Second * 5)
	var stop = ctrl.SetupSignalHandler()

	for {
		select {
		case <-tick.C:
			holdingRegisters := server.HoldingRegisters
			holdingRegisters[0] = uint16(rand.Intn(400))
			logrus.Info("Changing holding register 0 to: ", holdingRegisters[0])

			holdingRegisters[1] = uint16(rand.Intn(400))
			logrus.Info("Changing holding register 1 to: ", holdingRegisters[1])

			if holdingRegisters[0] > holdingRegisters[5] || holdingRegisters[1] > holdingRegisters[5] {
				server.Coils[0] = 1
			} else {
				server.Coils[0] = 0
			}
		case <-stop:
			return
		}
	}
}
