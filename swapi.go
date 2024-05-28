package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache search terms and people cache
var cacheMap = map[string]*cache.Cache{
	"people":    cache.New(1*time.Hour, 1*time.Hour),
	"species":   cache.New(1*time.Hour, 1*time.Hour),
	"planets":   cache.New(1*time.Hour, 1*time.Hour),
	"vehicles":  cache.New(1*time.Hour, 1*time.Hour),
	"starships": cache.New(1*time.Hour, 1*time.Hour),
	"films":     cache.New(1*time.Hour, 1*time.Hour),
}
var peopleCache = *cache.New(1*time.Hour, 1*time.Hour)

// Define Swapi Api here so we can mock for testing
type SwapiApi struct {
	handler func(url string) (*http.Response, error)
}

func (s SwapiApi) Get(url string) (*http.Response, error) {
	return s.handler(url)
}

func CachedSearchEntities[T any](path string, searchTerm string, swapiApi SwapiApi) ([]*T, error) {
	if entities, found := cacheMap[path].Get(searchTerm); found {
		return entities.([]*T), nil
	}

	entities, err := SearchEntities[T](path, searchTerm, swapiApi)
	if err != nil {
		return nil, err
	}

	cacheMap[path].Set(searchTerm, entities, cache.DefaultExpiration)
	return entities, nil
}

// Search for entities by search term
func SearchEntities[T any](path string, searchTerm string, swapiApi SwapiApi) ([]*T, error) {
	type Response struct {
		Count    int    `json:"count"`
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []*T   `json:"results"`
	}
	entities := []*T{}

	next := fmt.Sprintf("https://swapi.dev/api/%s/?search=%s&page=1", path, searchTerm)

	for next != "" {
		resp, err := swapiApi.Get(next)
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

// Fetch a person by URL
func CachedFetchPerson(url string, swapiApi SwapiApi) (*Person, error) {
	if person, found := peopleCache.Get(url); found {
		return person.(*Person), nil
	}

	resp, err := swapiApi.Get(url)
	if err != nil {
		fmt.Println("Error sending people request")
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading people response")
		return nil, err
	}

	var person Person
	err = json.Unmarshal(body, &person)
	if err != nil {
		fmt.Println("Error unmarshalling people response")
		return nil, err
	}

	peopleCache.Set(url, &person, cache.DefaultExpiration)

	return &person, nil
}

// Pre-fetch all people from SWAPI. This is a one-time operation.
func PreFetchPeople() error {
	next := "https://swapi.dev/api/people/"

	for next != "" {
		resp, err := http.Get(next)
		if err != nil {
			fmt.Println("Error sending people request")
			return err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading people response")
			return err
		}

		var response struct {
			Count    int       `json:"count"`
			Next     string    `json:"next"`
			Previous string    `json:"previous"`
			Results  []*Person `json:"results"`
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Error unmarshalling people response")
			return err
		}

		// Add all people to the cache
		for _, person := range response.Results {
			peopleCache.Set(person.URL, person, cache.DefaultExpiration)
		}

		next = response.Next
	}
	return nil
}
