package main

import (
	"log"
	"os"
	"time"

	"github.com/goburrow/modbus"
)

func main() {
	testTCP()
	//testRTU()
}
func testTCP() {
	// Modbus TCP
	handler := modbus.NewTCPClientHandler("127.0.0.1:5020")
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 0x01
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	err := handler.Connect()
	defer handler.Close()
	if err != nil {
		log.Fatal("fail to connect", err)
	}

	client := modbus.NewClient(handler)

	//results, err := client.ReadHoldingRegisters(1, 1)
	//if err != nil {
	//	log.Fatal("fail to read", err)
	//}

	//log.Println("result", results)
	results, err := client.ReadCoils(0, 1)

	log.Println("result", results)
}

func testRTU() {
	// Modbus RTU/ASCII
	handler := modbus.NewRTUClientHandler("/dev/ttys007")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second

	err := handler.Connect()
	if err != nil {
		log.Fatal("error connecting: ", err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadHoldingRegisters(0, 1)
	if err != nil {
		log.Fatal("error reading: ", err)
	}
	log.Println(results)
}
