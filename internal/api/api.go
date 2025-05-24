package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonResponse interface {
	GetResults() []string
}

// https://pokeapi.co/api/v2/location-area/{id or name}/
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

func GetLocationArea(url string) (response LocationResponse, err error) {
	response = LocationResponse{}
	err = getResponse(url, &response)
	if err != nil {
		fmt.Println(err)
		return LocationResponse{}, err
	}
	return response, nil
}

func getResponse(url string, response JsonResponse) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP request failed with %v | status: %v", url, res.Status)
		return err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return err
	}

	return nil
}
