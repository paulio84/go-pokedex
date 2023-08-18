package pokeapi

import (
	"encoding/json"
	"errors"
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
