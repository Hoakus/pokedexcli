package api

import (
	"encoding/json"
	"fmt"
	"github.com/Hoakus/pokedexcli/internal/pokecache"
	"io"
	"net/http"
	"time"
)

type JsonResponse interface{}

type LocationResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (l LocationResponse) GetResults() []string {
	var names []string
	for _, v := range l.Results {
		names = append(names, v.Name)
	}
	return names
}

type ExploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (e ExploreResponse) GetResults() []string {
	var names []string
	for _, v := range e.PokemonEncounters {
		names = append(names, v.Pokemon.Name)
	}
	return names
}

type PokemonResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	BaseXp    int    `json:"base_experience"`
	Height    int    `json:"height"`
	Weight    int    `json:"weight"`
	IsDefault bool   `json:"is_default"`
	Order     int    `json:"order"`
	Stats     []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func GetLocationArea(url string) (response LocationResponse, err error) {
	response = LocationResponse{}
	err = getResponse(url, &response)
	if err != nil {
		fmt.Println(err)
		return LocationResponse{}, err
	}
	return response, nil
}

func GetAreaByName(url string) (response ExploreResponse, err error) {
	response = ExploreResponse{}
	err = getResponse(url, &response)
	if err != nil {
		fmt.Println(err)
		return ExploreResponse{}, err
	}
	return response, nil
}

func GetPokemonByName(url string) (response PokemonResponse, err error) {
	response = PokemonResponse{}
	err = getResponse(url, &response)
	if err != nil {
		fmt.Println(err)
		return PokemonResponse{}, err
	}
	return response, nil
}

var pokeCache *pokecache.Cache

func getResponse(url string, response JsonResponse) error {
	if pokeCache == nil {
		pokeCache = pokecache.NewCache(time.Second * 30)
	}

	if val, ok := pokeCache.Get(url); ok {
		err := json.Unmarshal(val, &response)
		if err != nil {
			return err
		}
		success := fmt.Sprintf("accessed cache | %v found", url)
		fmt.Println(success)
		return nil
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP request failed with %v | status: %v", url, res.Status)
		return err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading results boddy")
		return err
	}

	json.Unmarshal(data, &response)
	pokeCache.Add(url, data)
	return nil
}
