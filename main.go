package main

import (
	"bufio"
	"fmt"
	api "internal/pokeapi"
	"os"
)

var commands map[string]cliCommand = make(map[string]cliCommand)

type cliCommand struct {
	name        string
	description string
	callback    func(*api.Config) error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	buildOptions()
	config := api.NewConfig()

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()

		if command, ok := commands[scanner.Text()]; ok {
			err := command.callback(config)
			if err != nil {
				fmt.Printf("error: %s\n", err.Error())
			}
		}
	}
}

func buildOptions() {
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	commands["map"] = cliCommand{
		name:        "map",
		description: "Explore the next 20 location areas in the Pokemon world",
		callback:    commandMap,
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Explore the previous 20 location areas in the Pokemon world",
		callback:    commandMapB,
	}
}

func commandMap(config *api.Config) error {
	err := config.Map()
	if err != nil {
		return err
	}

	for _, result := range config.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapB(config *api.Config) error {
	err := config.MapB()
	if err != nil {
		return err
	}

	for _, result := range config.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandHelp(config *api.Config) error {
	fmt.Println()
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage: \n")
	fmt.Println()

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()

	return nil
}

func commandExit(config *api.Config) error {
	os.Exit(0)
	return nil
}
