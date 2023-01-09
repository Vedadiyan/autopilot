package main

import (
	"fmt"
	"os"

	flaggy "github.com/vedadiyan/flaggy/pkg"
)

func main() {
	options := Options{}
	err := flaggy.Parse(&options, os.Args[1:])
	if err != nil {
		fmt.Println(err.Error())
	}
}
