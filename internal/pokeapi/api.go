package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/paulio84/go-pokedex/internal/pokecache"
)

var cache *pokecache.Cache

func NewAPI() *Api {
	defaultNext := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	cache = pokecache.NewCache(1 * time.Minute)

	return &Api{
		mapConfig: &mapConfig{
			Next:     &defaultNext,
			Previous: nil,
			Results:  nil,
		},
		explore: &explore{
			Name:              nil,
			PokemonEncounters: nil,
		},
	}
}

func (api *Api) Map() ([]Result, error) {
	err := api.mapConfig.mapF()

	return api.mapConfig.convertMapResults(), err
}

func (c *mapConfig) mapF() error {
	if c.Next == nil {
		return errors.New("at the edge of the Pokemon world")
	}

	err := c.getLocations(*c.Next)
	if err != nil {
		return err
	}

	return nil
}

func (api *Api) MapB() ([]Result, error) {
	err := api.mapConfig.mapB()

	return api.mapConfig.convertMapResults(), err
}

func (c *mapConfig) mapB() error {
	if c.Previous == nil {
		return errors.New("cannot go back any further")
	}

	err := c.getLocations(*c.Previous)
	if err != nil {
		return err
	}

	return nil
}

func (c *mapConfig) getLocations(url string) error {
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

func (api *Api) Explore(areaName string) ([]Result, error) {
	err := api.explore.explore(areaName)

	return api.explore.convertPokemonEncounters(), err
}

func (e *explore) explore(areaName string) error {
	err := e.getPokemonByArea(areaName)
	if err != nil {
		return err
	}

	return nil
}

func (e *explore) getPokemonByArea(areaName string) error {
	var responseBody []byte
	var ok bool
	var err error

	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	// try to get the response body from the cache
	responseBody, ok = cache.Get(url)
	if !ok {
		// get the response from the network
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 404 {
			return fmt.Errorf("invalid area name: %s", areaName)
		}

		responseBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// store the response in the cache
		cache.Add(url, responseBody)
	}

	err = json.Unmarshal(responseBody, &e)
	if err != nil {
		return err
	}

	return nil
}
