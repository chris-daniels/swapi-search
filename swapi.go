package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchEntities[T any](path string) ([]*T, error) {
	type Response struct {
		Count    int    `json:"count"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []*T   `json:"results"`
	}
	entities := []*T{}

	next := fmt.Sprintf("https://swapi.dev/api/%s/?page=1", path)

	for next != "" {
		resp, err := http.Get(next)
		if err != nil {
			fmt.Println("Error sending people request")
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading people response")
			return nil, err
		}

		var response Response
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Error unmarshalling people response")
			return nil, err
		}

		entities = append(entities, response.Results...)
		next = response.Next
	}

	return entities, nil
}
