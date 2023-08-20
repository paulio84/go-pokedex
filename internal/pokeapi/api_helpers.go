package pokeapi

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/paulio84/go-pokedex/internal/pokecache"
)

var cache *pokecache.Cache

func initCache() *pokecache.Cache {
	return pokecache.NewCache(1 * time.Minute)
}

func getLocations(url string) ([]byte, error) {
	responseBody, _, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func getPokemonByArea(areaName string) ([]byte, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	responseBody, statusCode, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 {
		return nil, fmt.Errorf("invalid area name: %s", areaName)
	}

	return responseBody, nil
}

func catchPokemon(pokemonName string) ([]byte, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	responseBody, statusCode, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 {
		return nil, fmt.Errorf("pokemon %s, does not exist", pokemonName)
	}

	return responseBody, nil
}

func makeRequest(url string) ([]byte, int, error) {
	// try to get the response body from the cache
	responseBody, ok := cache.Get(url)
	if !ok {
		// get the response from the network
		resp, err := http.Get(url)
		if err != nil {
			return nil, resp.StatusCode, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 404 {
			return nil, resp.StatusCode, nil
		}

		responseBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, resp.StatusCode, err
		}

		// store the response in the cache
		cache.Add(url, responseBody)
	}

	return responseBody, 0, nil
}

func baseExpChance(baseExp int) int {
	switch {
	case baseExp <= 50:
		return 8
	case baseExp >= 51 && baseExp <= 100:
		return 7
	case baseExp >= 101 && baseExp <= 200:
		return 6
	case baseExp >= 201 && baseExp <= 300:
		return 5
	case baseExp >= 301 && baseExp <= 400:
		return 4
	case baseExp >= 401 && baseExp <= 500:
		return 3
	case baseExp >= 501 && baseExp <= 600:
		return 2
	default:
		return 1
	}
}

func newMapConfig() *mapConfig {
	defaultNext := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"

	return &mapConfig{
		Next:     &defaultNext,
		Previous: nil,
		Results:  nil,
	}
}

func newExplore() *explore {
	return &explore{
		Name:              nil,
		PokemonEncounters: nil,
	}
}

func newPokemon() *Pokemon {
	return &Pokemon{
		Name:           nil,
		Height:         nil,
		Weight:         nil,
		BaseExperience: nil,
		Stats:          nil,
		Types:          nil,
	}
}
