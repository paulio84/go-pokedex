package main

import (
	"bufio"
	"fmt"
	"os"
)

var options map[string]cliCommand = make(map[string]cliCommand)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	buildOptions()

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()

		if command, ok := options[scanner.Text()]; ok {
			err := command.callback()
			if err != nil {
				fmt.Printf("error: %s", err.Error())
				break
			}
		}
	}
}

func buildOptions() {
	options["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}
	options["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
}

func commandHelp() error {
	fmt.Printf("\nWelcome to the Pokedex!\n")
	fmt.Printf("Usage: \n")
	fmt.Println()

	for _, command := range options {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()

	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}
