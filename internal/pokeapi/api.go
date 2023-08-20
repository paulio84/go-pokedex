package pokeapi

import (
	"encoding/json"
	"errors"
	"math/rand"
)

func NewAPI() *Api {
	cache = initCache()

	return &Api{
		mapConfig: newMapConfig(),
		explore:   newExplore(),
		pokemon:   newPokemon(),
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

	responseBody, err := getLocations(*c.Next)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &c)
	return err
}

func (api *Api) MapB() ([]Result, error) {
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

func (api *Api) Explore(areaName string) ([]Result, error) {
	err := api.explore.explore(areaName)

	return api.explore.convertPokemonEncounters(), err
}

func (e *explore) explore(areaName string) error {
	responseBody, err := getPokemonByArea(areaName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &e)
	return err
}

func (api *Api) Catch(pokemonName string) (Pokemon, error, bool) {
	err := api.pokemon.catch(pokemonName)
	var caught bool

	// determine whether the pokemon is caught
	if err == nil {
		caught = true
		chance := baseExpChance(*api.pokemon.BaseExperience)
		roll := rand.Intn(10) + 1 // roll a 1-10
		if roll > chance {
			api.pokemon = newPokemon()
			caught = false
		}
	}

	return *api.pokemon, err, caught
}

func (p *Pokemon) catch(pokemonName string) error {
	responseBody, err := catchPokemon(pokemonName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &p)
	return err
}
