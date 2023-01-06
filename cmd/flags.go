package main

import (
	"fmt"
	"os"

	"github.com/vedadiyan/autopilot/internal/microservices"
	flaggy "github.com/vedadiyan/flaggy/pkg"
	getup "github.com/vedadiyan/getup/pkg"
)

type Options struct {
	Microservice microservices.Microservice `long:"microservice" short:"" help:"Automatically generates a microservice"`
	Setup        bool                       `long:"setup" short:"" help:"Setups autopilot in the system"`
}

func (o Options) Run() error {
	if o.Setup {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		getup.New(homeDir, "gopher", "autopilot")
		err = getup.Setup()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		fmt.Println("autopilot successfully setup in the system")
		return nil
	}
	flaggy.PrintHelp()
	return nil
}
