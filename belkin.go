package main

import (
	"github.com/gorilla/mux"
	"github.com/savaki/wemo"
	"io"
	"net/http"
)

type BelkinRequest struct {
	name   string
	action string
}

func BelkinProcessor(ch chan BelkinRequest) {
	for {
		request := <-ch

		var devices []*wemo.Device
		switch request.name {
		case "office_plug":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.32:49153"},
			}

		case "kitchen_overhead":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.8:49154"},
			}

		case "kitchen_sink":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.18:49153"},
			}

		case "mirror_light":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.17:49153"},
			}

		case "office_overhead":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.19:49153"},
			}

		case "all":
			devices = []*wemo.Device{
				&wemo.Device{Host: "10.0.1.32:49153"},
				&wemo.Device{Host: "10.0.1.8:49154"},
				&wemo.Device{Host: "10.0.1.18:49153"},
				&wemo.Device{Host: "10.0.1.17:49153"},
				&wemo.Device{Host: "10.0.1.19:49153"},
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
			}
		}
	}
}

func BelkinHandler(ch chan BelkinRequest) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		name := vars["name"]
		action := vars["action"]

		ch <- BelkinRequest{name: name, action: action}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "hello world")
	}
}
