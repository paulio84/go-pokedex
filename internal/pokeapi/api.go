package pokeapi

import (
	"encoding/json"
	"errors"
	"math/rand"
)

func NewAPI() *API {
	cache = initialiseCache()

	return &API{
		mapConfig: newMapConfig(),
	}
}

func (api *API) Map() ([]Result, error) {
	err := api.mapConfig.mapF()

	return api.mapConfig.convertMapResults(), err
}

func (c *mapConfig) mapF() error {
	if c.Next == nil {
		return errors.New("at the edge of the Pokemon world")
	}

	responseBody, err := getLocations(*c.Next)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &c)
	return err
}

func (api *API) MapB() ([]Result, error) {
	err := api.mapConfig.mapB()

	return api.mapConfig.convertMapResults(), err
}

func (c *mapConfig) mapB() error {
	if c.Previous == nil {
		return errors.New("cannot go back any further")
	}

	responseBody, err := getLocations(*c.Previous)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &c)
	return err
}

func (api *API) Explore(areaName string) ([]Result, error) {
	exploreArea, err := exploreArea(areaName)

	return exploreArea.convertPokemonEncounters(), err
}

func exploreArea(areaName string) (explore, error) {
	var explore explore

	responseBody, err := getPokemonByArea(areaName)
	if err != nil {
		return explore, err
	}

	err = json.Unmarshal(responseBody, &explore)
	return explore, err
}

func (api *API) Catch(pokemonName string) (Pokemon, bool, error) {
	pokemon, err := catch(pokemonName)
	var caught bool

	// determine whether the pokemon is caught
	if err == nil {
		chance := baseExpChance(*pokemon.BaseExperience)
		roll := rand.Intn(10) + 1 // roll a 1-10
		if roll <= chance {
			caught = true
		}
	}

	return pokemon, caught, err
}

func catch(pokemonName string) (Pokemon, error) {
	var pokemon Pokemon

	responseBody, err := catchPokemon(pokemonName)
	if err != nil {
		return pokemon, err
	}

	err = json.Unmarshal(responseBody, &pokemon)
	return pokemon, err
}
