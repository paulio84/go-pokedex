package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/paulio84/go-pokedex/internal/pokeapi"
)

var commands map[string]cliCommand
var api *pokeapi.Api
var pokedex map[string]pokeapi.Pokemon

type cliCommand struct {
	name        string
	description string
	callback    func(string) error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	api = pokeapi.NewAPI()
	buildCommands()
	pokedex = make(map[string]pokeapi.Pokemon)

	for {
		arg := ""

		fmt.Print("pokedex > ")
		scanner.Scan()

		words := strings.Fields(scanner.Text())
		if len(words) > 1 {
			arg = words[1]
		}

		if command, ok := commands[words[0]]; ok {
			err := command.callback(arg)
			if err != nil {
				fmt.Printf("error: %s\n", err.Error())
			}
		}
	}
}

func buildCommands() {
	commands = make(map[string]cliCommand)

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
		description: "Display the next location areas in the Pokemon world",
		callback:    commandMap,
	}
	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Display the previous location areas in the Pokemon world",
		callback:    commandMapB,
	}
	commands["explore"] = cliCommand{
		name:        "explore [area_name]",
		description: "Explore a given area [area_name] and display the Pokemon found there",
		callback:    commandExplore,
	}
	commands["catch"] = cliCommand{
		name:        "catch [pokemon_name]",
		description: "Attempt to catch the Pokemon [pokemon_name]",
		callback:    commandCatch,
	}
}

func commandCatch(arg string) error {
	if arg == "" {
		return errors.New("you must enter a pokemon name")
	}

	if _, ok := pokedex[arg]; ok {
		return fmt.Errorf("you've already caught %s", arg)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", arg)

	pokemon, err := api.Catch(arg)
	if err != nil {
		return err
	}

	if pokemon != nil {
		fmt.Printf("%s was caught!\n", arg)
		// add pokemon to our pokedex
		pokedex[arg] = *pokemon
		return nil
	}

	fmt.Printf("%s escaped!\n", arg)
	return nil
}

func commandExplore(arg string) error {
	if arg == "" {
		return errors.New("you must enter an area name to explore")
	}

	results, err := api.Explore(arg)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", arg)
	fmt.Println("Found Pokemon:")
	displayResults(results)

	return nil
}

func commandMap(arg string) error {
	results, err := api.Map()
	if err != nil {
		return err
	}

	displayResults(results)

	return nil
}

func commandMapB(arg string) error {
	results, err := api.MapB()
	if err != nil {
		return err
	}

	displayResults(results)

	return nil
}

func commandHelp(arg string) error {
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

func commandExit(arg string) error {
	os.Exit(0)
	return nil
}

func displayResults(results []pokeapi.Result) {
	for _, result := range results {
		fmt.Println(result.Display())
	}
}
