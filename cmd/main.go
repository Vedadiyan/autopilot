package main

import (
	"os"

	flaggy "github.com/vedadiyan/flaggy/pkg"
)

func main() {
	options := Options{}
	flaggy.Parse(&options, os.Args[1:])
}
