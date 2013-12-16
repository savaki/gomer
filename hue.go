package main

import (
	"github.com/gorilla/mux"
	"github.com/savaki/go.hue"
	"io"
	"net/http"
)

type HueRequest struct {
	name   string
	action string
}

func HueProcessor(username string, ch chan HueRequest) {
	bridge := hue.NewBridge("10.0.1.11", username)
	for {
		request := <-ch

		var lights []*hue.Light = nil
		switch request.name {
		case "all":
			lights, _ = bridge.GetAllLights()

		case "spotlight":
			spotlight, err := bridge.FindLightById("9")
			if err == nil {
				lights = append(lights, spotlight)
			}

		case "dome":
			left, err := bridge.FindLightById("3")
			if err == nil {
				lights = append(lights, left)
			}

			right, err := bridge.FindLightById("6")
			if err == nil {
				lights = append(lights, right)
			}

		case "front_room":
			door, err := bridge.FindLightById("10")
			if err == nil {
				lights = append(lights, door)
			}

			sofa, err := bridge.FindLightById("4")
			if err == nil {
				lights = append(lights, sofa)
			}
		}

		if lights != nil {
			for _, light := range lights {
				switch request.action {
				case "toggle":
					attributes, err := light.GetLightAttributes()
					if err == nil {
						if attributes.State.On {
							light.Off()
						} else {
							light.On()
						}
					}

				case "on":
					light.On()

				case "off":
					light.Off()

				case "colorloop":
					light.ColorLoop()
				}
			}
		}
	}
}

func HueHandler(ch chan HueRequest) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		name := vars["name"]
		action := vars["action"]

		ch <- HueRequest{name: name, action: action}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "hello world")
	}
}
