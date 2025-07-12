package main

import (
	"fmt"
	"os"
)

func callbackExit(cleanup func()) func() error {
	return func() error {
		if cleanup != nil {
			cleanup()
		}
		fmt.Println("\033[1;31m Exiting...\033[0m")

		os.Exit(0)
		return nil
	}
}
