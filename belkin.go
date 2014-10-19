package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/savaki/wemo"
	"log"
	"strconv"
)

type BelkinRequest struct {
	name     string
	action   string
	response chan string
}

func BelkinProcessor(ch chan BelkinRequest, deviceConfigs map[string]string) {
	for {
		request := <-ch

		var devices []*wemo.Device
		if host, found := deviceConfigs[request.name]; found {
			devices = []*wemo.Device{
				&wemo.Device{Host: host},
			}
		} else if request.name == "all" {
			devices = []*wemo.Device{}
			for _, host := range deviceConfigs {
				devices = append(devices, &wemo.Device{
					Host: host,
				})
			}
		} else {
			log.Printf("WARNING - unknown device name, %s\n", request.name)
		}

		for _, device := range devices {
			switch request.action {
			case "on":
				device.On()

			case "off":
				device.Off()

			case "toggle":
				device.Toggle()

			case "state":
				if request.response != nil {
					result := strconv.Itoa(device.GetBinaryState())
					request.response <- result
				}
			}
		}
	}
}

func BelkinHandler(ch chan BelkinRequest) func(*gin.Context) {
	return func(c *gin.Context) {
		name := c.Params.ByName("name")
		action := c.Params.ByName("action")

		var response chan string = nil
		if action == "state" {
			response = make(chan string)
		}

		ch <- BelkinRequest{name: name, action: action, response: response}

		if response != nil {
			defer close(response)
			value := <-response
			json := fmt.Sprintf(`{"status":"ok","state":%s}`, value)
			c.JSON(200, json)

		} else {
			c.JSON(200, map[string]string{"status": "ok"})
		}
	}
}
