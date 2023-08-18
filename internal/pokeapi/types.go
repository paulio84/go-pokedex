package pokeapi

type Result interface {
	Display() string
}

type Api struct {
	mapConfig *mapConfig
}

type mapConfig struct {
	Next     *string
	Previous *string
	Results  []mapResult
}

func (mc mapConfig) convertMapResults() []Result {
	results := make([]Result, len(mc.Results))
	for i, v := range mc.Results {
		results[i] = Result(v)
	}
	return results
}

type mapResult struct {
	Name string
}

func (mr mapResult) Display() string {
	return mr.Name
}
