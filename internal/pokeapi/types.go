package pokeapi

import "fmt"

type Pokemon struct {
	Name           *string        `json:"name"`
	Height         *int           `json:"height"`
	Weight         *int           `json:"weight"`
	BaseExperience *int           `json:"base_experience"`
	Stats          []pokemonStats `json:"stats"`
	Types          []pokemonTypes `json:"types"`
}

type pokemonStats struct {
	BaseStat int         `json:"base_stat"`
	Stat     pokemonStat `json:"stat"`
}

type pokemonStat struct {
	Name string `json:"name"`
}

type pokemonTypes struct {
	Type pokemonType `json:"type"`
}

type pokemonType struct {
	Name string `json:"name"`
}

type Result interface {
	Display() string
}

type Api struct {
	mapConfig *mapConfig
	explore   *explore
	pokemon   *Pokemon
}

type mapConfig struct {
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
	Results  []mapResult `json:"results"`
}

func (mc mapConfig) convertMapResults() []Result {
	results := make([]Result, len(mc.Results))
	for i, v := range mc.Results {
		results[i] = Result(v)
	}
	return results
}

type mapResult struct {
	Name string `json:"name"`
}

func (mr mapResult) Display() string {
	return mr.Name
}

type explore struct {
	Name              *string             `json:"name"`
	PokemonEncounters []pokemonEncounters `json:"pokemon_encounters"`
}

func (e explore) convertPokemonEncounters() []Result {
	results := make([]Result, len(e.PokemonEncounters))
	for i, v := range e.PokemonEncounters {
		results[i] = Result(v)
	}
	return results
}

type pokemon struct {
	Name string `json:"name"`
	// URL  string `json:"url"`
}

type pokemonEncounters struct {
	Pokemon pokemon `json:"pokemon"`
}

func (pe pokemonEncounters) Display() string {
	return fmt.Sprintf(" - %s", pe.Pokemon.Name)
}
