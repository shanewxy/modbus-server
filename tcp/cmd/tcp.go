package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tbrandon/mbserver"
	ctrl "sigs.k8s.io/controller-runtime"
)

var device TCPDevice

func NewCommand() *cobra.Command {
	var c = &cobra.Command{
		Use:  "tcp",
		Long: "test device tcp",
		RunE: func(cmd *cobra.Command, args []string) error {
			device.Run()
			return nil
		},
	}

	c.Flags().StringVarP(&device.ip, "ip", "i", "0.0.0.0", "tcp listening ip address of the modbus server")
	c.Flags().IntVarP(&device.port, "port", "p", 5020, "tcp listening port of the modbus server")

	return c
}

type TCPDevice struct {
	ip   string
	port int
}

func (r *TCPDevice) Run() {
	endpoint := fmt.Sprintf("%s:%d", r.ip, r.port)
	server := mbserver.NewServer()
	if err := server.ListenTCP(endpoint); err != nil {
		logrus.Error("Fail to start server", err)
		return
	}
	logrus.Info("Listening on ", endpoint)

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
				logrus.Info("Alerting")
				server.Coils[0] = 1
			} else {
				server.Coils[0] = 0
			}
		case <-stop:
			return
		}
	}
}
