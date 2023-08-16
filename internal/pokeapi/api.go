package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/paulio84/go-pokedex/internal/pokecache"
)

type Config struct {
	Next     *string
	Previous *string
	Results  []result
}

type result struct {
	Name string
}

var cache *pokecache.Cache

func NewConfig() *Config {
	defaultNext := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	cache = pokecache.NewCache(10 * time.Second)

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

	err := c.getLocations(*c.Next)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) MapB() error {
	if c.Previous == nil {
		return errors.New("cannot go back any further")
	}

	err := c.getLocations(*c.Previous)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) getLocations(url string) error {
	var responseBody []byte
	var ok bool
	var err error

	// try to get the response body from the cache
	responseBody, ok = cache.Get(url)
	if !ok {
		// get the response from the network
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		responseBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// store the response in the cache
		cache.Add(url, responseBody)
	}

	err = json.Unmarshal(responseBody, &c)
	if err != nil {
		return err
	}

	return nil
}
