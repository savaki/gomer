package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gomer"
	app.Usage = "home automation system"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.IntFlag{"port", 8080, "the port number to list on", "PORT"},
		cli.StringFlag{"hue", "", "the hue username", "HUE_USERNAME"},
		cli.StringFlag{"docroot", "public", "html content directory", "DOCROOT"},
	}
	app.Action = Run
	app.Run(os.Args)
}

func Run(c *cli.Context) {
	port := c.Int("port")
	username := c.String("hue")
	docroot := c.String("docroot")

	routes := gin.New()

	// handle belkin related commands
	belkinCh := make(chan BelkinRequest, 20)
	go BelkinProcessor(belkinCh)
	routes.POST("/api/belkin/:name/:action", BelkinHandler(belkinCh))

	// handle hue related commands
	hueCh := make(chan HueRequest, 20)
	go HueProcessor(username, hueCh)
	routes.POST("/api/hue/:name/:action", HueHandler(hueCh))

	// handle ir related commands
	irCh := make(chan IrRequest, 20)
	go IrProcessor(irCh)
	routes.POST("/api/ir/:name/:action", IrHandler(irCh))

	// default to file server
	routes.Static("/", docroot)

	http.ListenAndServe(fmt.Sprintf(":%d", port), routes)
}
