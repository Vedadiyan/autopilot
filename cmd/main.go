package main

import (
	"autopilot/internal/ssh"
	"fmt"
	"runtime"
)

func main() {
	ssh := ssh.New("45.77.39.116", 22, "root", "******")
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
	if err := ssh.ExecMany([]string{
		"podman rm --force Redis-01",
		"podman ps",
		"pwd",
	}, output, errors); err != nil {
		panic(err)
	}
	running = false
	fmt.Println("Completed")
	runtime.Goexit()
}
