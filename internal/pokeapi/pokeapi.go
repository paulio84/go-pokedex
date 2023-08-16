package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Config struct {
	Next     *string
	Previous *string
	Results  []Result
}

type Result struct {
	Name string
}

func NewConfig() (c *Config) {
	defaultNext := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"

	return &Config{
		Next:     &defaultNext,
		Previous: nil,
		Results:  nil,
	}
}

func (c *Config) Map() error {
	if c.Next == nil {
		return errors.New("at the end of the Pokemon world")
	}

	responseBody, err := getLocations(*c.Next)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) MapB() error {
	if c.Previous == nil {
		return errors.New("cannot go back any further")
	}

	responseBody, err := getLocations(*c.Previous)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &c)
	if err != nil {
		return err
	}

	return nil
}

func getLocations(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
