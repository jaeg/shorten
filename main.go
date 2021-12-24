package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/jaeg/shorten/app"
)

func main() {
	app := &app.App{}
	app.Init()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c

		cancel()
	}()

	app.Run(ctx)
}
