package main

import (
	"fmt"
	"growibackup"
	"os"
)

func main() {
	err := growibackup.Run(os.Args[1], os.Args[2])
	exitCode := 0
	if err != nil {
		fmt.Println(err)
		exitCode = 1
	}
	os.Exit(exitCode)
}
