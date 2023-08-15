package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var options map[string]cliCommand = make(map[string]cliCommand)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     *string
	Previous *string
	Results  []Result
}

func (c *Config) updateNext(v string) {
	c.Next = &v
}

type Result struct {
	Name string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	buildOptions()
	config := &Config{Next: nil, Previous: nil, Results: nil}

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()

		if command, ok := options[scanner.Text()]; ok {
			err := command.callback(config)
			if err != nil {
				fmt.Printf("error: %s\n", err.Error())
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
	options["map"] = cliCommand{
		name:        "map",
		description: "Explore the next 20 location areas in the Pokemon world",
		callback:    commandMap,
	}
	options["mapb"] = cliCommand{
		name:        "mapb",
		description: "Explore the previous 20 location areas in the Pokemon world",
		callback:    commandMapB,
	}
}

func commandMap(config *Config) error {
	if config.Next == nil {
		config.updateNext("https://pokeapi.co/api/v2/location-area?offset=0&limit=20")
	}

	// TODO: Move this API call to an internal package
	resp, err := http.Get(*config.Next)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &config)
	if err != nil {
		return err
	}

	for _, result := range config.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapB(config *Config) error {
	if config.Previous == nil {
		return errors.New("cannot go back any further")
	}

	// TODO: Move this API call to an internal package
	resp, err := http.Get(*config.Previous)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &config)
	if err != nil {
		return err
	}

	for _, result := range config.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandHelp(config *Config) error {
	fmt.Println()
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Printf("Usage: \n")
	fmt.Println()

	for _, command := range options {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println()

	return nil
}

func commandExit(config *Config) error {
	os.Exit(0)
	return nil
}
