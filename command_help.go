package main

import "fmt"

func callbackHelp() error {
	fmt.Println("\nHere are your available commands:")

	availableCommands := getCommands(nil)

	for _, cmd := range availableCommands {
		fmt.Printf(" - %s: %s\n", cmd.name, cmd.description)
	}

	fmt.Println("")

	return nil
}
