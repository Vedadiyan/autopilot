package main

import (
	"autopilot/internal/ssh"
	"fmt"
	"runtime"
)

func main() {
	ssh := ssh.New("45.77.39.116", 22, "root", "Pou9617!")
	if err := ssh.Connect(); err != nil {
		panic(err)
	}
	output := make(chan string)
	errors := make(chan string)
	running := true
	go func() {
		for running {
			select {
			case out := <-output:
				{
					fmt.Println(out)
				}
			case err := <-errors:
				{
					fmt.Println(err)
				}
			}
		}
	}()
	if err := ssh.Exec("dnf install -y git", output, errors); err != nil {
		panic(err)
	}
	running = false
	fmt.Println("Completed")
	runtime.Goexit()
}
