package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/savaki/wemo"
	"io"
	"net/http"
	"strconv"
)

type BelkinRequest struct {
	name     string
	action   string
	response chan string
}

func BelkinProcessor(ch chan BelkinRequest) {
	for {
		request := <-ch

		var devices []*wemo.Device
		switch request.name {
		case "kitchen_overhead":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.8:49154"},
			}

		case "kitchen_sink":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.18:49153"},
			}

		case "mirror_overhead":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.17:49153"},
			}

		case "office_overhead":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.19:49154"},
				// &wemo.Device{Host: "10.0.1.19:49153"},
			}

		case "bathroom":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.24:49153"},
			}

		case "left_wall":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.32:49153"},
			}

		case "right_wall":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.25:49153"},
			}

		case "all":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.32:49153"},
				&wemo.Device{Host: "10.0.1.8:49154"},
				&wemo.Device{Host: "10.0.1.18:49153"},
				&wemo.Device{Host: "10.0.1.17:49153"},
				&wemo.Device{Host: "10.0.1.19:49154"},
				// &wemo.Device{Host: "10.0.1.19:49153"},
				&wemo.Device{Host: "10.0.1.25:49153"},
			}
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

func BelkinHandler(ch chan BelkinRequest) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		name := vars["name"]
		action := vars["action"]

		var response chan string = nil
		if action == "state" {
			response = make(chan string)
		}

		ch <- BelkinRequest{name: name, action: action, response: response}

		w.WriteHeader(http.StatusOK)

		if response != nil {
			defer close(response)
			value := <-response
			json := fmt.Sprintf(`{"status":"ok","state":%s}`, value)
			io.WriteString(w, json)

		} else {
			io.WriteString(w, `{"status":"ok"}`)
		}
	}
}
