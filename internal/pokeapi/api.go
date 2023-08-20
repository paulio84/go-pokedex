package pokeapi

import (
	"encoding/json"
	"errors"
	"math/rand"
)

func NewAPI() *Api {
	defaultNext := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	cache = initCache()

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
		pokemon: &Pokemon{
			Name:           nil,
			Height:         nil,
			Weight:         nil,
			BaseExperience: nil,
			Stats:          nil,
			Types:          nil,
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

func (api *Api) Catch(pokemonName string) (*Pokemon, error) {
	err := api.pokemon.catch(pokemonName)

	// determine whether the pokemon is caught
	chance := baseExpChance(*api.pokemon.BaseExperience)
	roll := rand.Intn(10) + 1 // roll a 1-10
	if roll > chance {
		api.pokemon = nil
	}

	return api.pokemon, err
}

func (p *Pokemon) catch(pokemonName string) error {
	responseBody, err := catchPokemon(pokemonName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &p)
	return err
}
