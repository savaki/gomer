package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"time"
)

func timeoutDialer(connectTimeout time.Duration, readWriteTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, connectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(readWriteTimeout))
		return conn, nil
	}
}

type IrRequest string

func IrProcessor(ch chan IrRequest) {
	for {
		request := <-ch

		// 1. create the ir command
		err := sendIr(request)
		if err != nil {
			fmt.Printf("send ir => %s\n", err.Error())
		}
	}
}

func toCommand(request IrRequest) string {
	var command string
	switch string(request) {
	case "projector_power":
		command = `sendir,1:3,5,37993,1,1,340,169,22,63,22,63,22,20,22,20,22,20,22,20,22,20,22,63,22,63,22,20,22,63,22,20,22,63,22,20,22,63,22,20,22,20,22,20,22,20,22,20,22,63,22,20,22,20,22,63,22,63,22,63,22,63,22,63,22,20,22,63,22,63,22,20,22,4863`

	case "itv_play_pause":
		command = `sendir,1:3,7,38461,1,1,347,173,21,22,21,65,21,65,21,65,21,22,21,65,21,65,21,65,21,65,21,65,21,65,21,22,21,22,21,22,21,22,21,65,21,65,21,65,21,65,21,65,21,65,21,22,21,65,21,22,21,22,21,65,21,65,21,22,21,65,21,22,21,65,21,65,21,1306,347,173,21,22,21,65,21,65,21,65,21,22,21,65,21,65,21,65,21,65,21,65,21,65,21,22,21,22,21,22,21,22,21,65,21,65,21,22,21,65,21,22,21,22,21,22,21,22,21,22,21,22,21,65,21,65,21,22,21,65,21,22,21,65,21,65,21,4923`

	case "itv_menu":
		command = `sendir,1:3,3,38461,1,1,347,172,21,22,21,64,21,64,21,64,21,22,21,64,21,64,21,64,21,64,21,64,21,64,21,22,21,22,21,22,21,22,21,64,21,64,21,64,21,22,21,22,21,22,21,22,21,22,21,22,21,22,21,64,21,64,21,22,21,64,21,22,21,64,21,65,21,1470,347,87,21,3695,347,87,21,4923`

	case "itv_left":
		command = `sendir,1:3,1,38580,1,1,348,173,21,22,21,65,21,65,21,65,21,22,21,65,21,65,21,65,21,65,21,65,21,65,21,22,21,22,21,22,21,22,21,65,21,65,21,22,21,22,21,65,21,22,21,22,21,22,21,22,21,22,21,65,21,65,21,22,21,65,21,22,21,65,21,65,21,1473,347,87,21,4938`

	case "itv_right":
		command = `sendir,1:3,1,38461,1,1,348,172,21,22,21,65,21,65,21,64,21,22,21,64,21,64,21,64,21,64,21,64,21,64,21,22,21,22,21,22,21,22,21,64,21,22,21,64,21,64,21,22,21,22,21,22,21,22,21,22,21,22,21,65,21,64,21,22,21,64,21,22,21,64,21,65,21,1470,347,87,21,4923`

	case "itv_up":
		command = `sendir,1:3,4,38580,1,1,347,173,21,22,21,65,21,65,21,65,21,22,21,65,21,65,21,65,21,65,21,65,21,65,21,22,21,22,21,22,21,22,21,65,21,22,21,65,21,22,21,65,21,22,21,22,21,22,21,22,21,22,21,65,21,65,21,22,21,65,21,22,21,65,21,65,21,1474,347,87,21,4938`

	case "itv_down":
		command = `sendir,1:3,2,38580,1,1,347,173,21,22,21,65,21,65,21,65,21,22,21,65,21,65,21,65,21,65,21,65,21,65,21,22,21,22,21,22,21,22,21,65,21,22,21,22,21,65,21,65,21,22,21,22,21,22,21,22,21,22,21,65,21,65,21,22,21,65,21,22,21,65,21,65,21,1475,347,87,21,4938`

	case "itv_select":
		command = `sendir,1:3,5,38461,1,1,347,173,21,22,21,65,21,65,21,65,21,22,21,65,21,65,21,65,21,65,21,65,21,65,21,22,21,22,21,22,21,22,21,65,21,22,21,22,21,65,21,65,21,65,21,22,21,65,21,22,21,22,21,65,21,65,21,22,21,65,21,22,21,65,21,65,21,1393,347,173,21,22,21,65,21,65,21,65,21,22,21,65,21,65,21,65,21,65,21,65,21,65,21,22,21,22,21,22,21,22,21,65,21,65,21,22,21,65,21,22,21,22,21,22,21,22,21,22,21,22,21,65,21,65,21,22,21,65,21,22,21,65,21,65,21,4923`

	case "yamaha_louder":
		command = `sendir,1:3,17,38226,1,1,343,172,22,21,22,64,22,21,22,64,22,64,22,64,22,64,22,21,22,64,22,21,22,64,22,21,22,21,22,21,22,21,22,64,22,21,22,64,22,21,22,64,22,64,22,21,22,21,22,21,22,64,22,21,22,64,22,21,22,21,22,64,22,64,22,64,22,1522,343,86,22,3669,343,86,22,3669,343,86,22,4892`

	case "yamaha_softer":
		command = `sendir,1:3,24,38226,1,1,343,172,22,21,22,64,22,21,22,64,22,64,22,64,22,64,22,21,22,64,22,21,22,64,22,21,22,21,22,21,22,21,22,64,22,64,22,64,22,21,22,64,22,64,22,21,22,21,22,21,22,21,22,21,22,64,22,21,22,21,22,64,22,64,22,64,22,1522,343,86,22,4892`

	}

	return command + "\r\n"
}

func sendIr(request IrRequest) error {
	command := toCommand(request)

	conn, err := timeoutDialer(2*time.Second, 2*time.Second)("tcp", "10.0.1.15:4998")
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("ir req => %s\n", request)
	_, err = conn.Write([]byte(command))
	if err != nil {
		fmt.Println("unable to write command")
		return err
	}

	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Printf("unable to read response => %d\n", n)
		return err
	}

	fmt.Printf("ir res => %s\n", string(data))
	return nil
}

func IrHandler(ch chan IrRequest) func(*gin.Context) {
	return func(c *gin.Context) {
		name := c.Params.ByName("name")

		ch <- IrRequest(name)

		c.String(200, "hello world")
	}
}
