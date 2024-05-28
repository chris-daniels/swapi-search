package main

// Write test for SearchValue function
// We should mock out requests to the SWAPI server

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestSearchValue(t *testing.T) {
	tests := []struct {
		name           string
		term           string
		searchResponse string
		peopleResponse string
		expected       string
	}{
		{
			name:           "Search for planets. Search response should return 1 planet with one person",
			term:           "tatooine",
			searchResponse: `{"count":1,"next":null,"previous":null,"results":[{"name":"Tatooine","residents":["https://swapi.dev/api/people/1/"]}]}`,
			peopleResponse: `{"name":"Luke Skywalker"}`,
			expected:       "Found 1 planet(s):\n\t - Tatooine\n\t\t - Luke Skywalker\n",
		},
		{
			name:           "Search for planets. Search response should return 1 planet with no people",
			term:           "hoth",
			searchResponse: `{"count":1,"next":null,"previous":null,"results":[{"name":"Hoth","residents":[]}]}`,
			peopleResponse: `{"name":"Luke Skywalker"}`,
			expected:       "Found 1 planet(s):\n\t - Hoth\n",
		},
		{
			name:           "Search for planets. Search response should return 0 planets that match search term",
			term:           "earth",
			searchResponse: `{"count":0,"next":null,"previous":null,"results":[]}`,
			peopleResponse: `{"name":"Luke Skywalker"}`,
			expected:       "",
		},
		{
			name:           "Search for planets. Search response should return 2 planets with 2 people",
			term:           "alder",
			searchResponse: `{"count":2,"next":null,"previous":null,"results":[{"name":"Alder","residents":["https://swapi.dev/api/people/1/","https://swapi.dev/api/people/2/"]},{"name":"Alderaan","residents":["https://swapi.dev/api/people/3/","https://swapi.dev/api/people/4/"]}]}`,
			peopleResponse: `{"name":"Luke Skywalker"}`,
			expected:       "Found 2 planet(s):\n\t - Alder\n\t\t - Luke Skywalker\n\t\t - Luke Skywalker\n\t - Alderaan\n\t\t - Luke Skywalker\n\t\t - Luke Skywalker\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSwapi := SwapiApi{
				handler: func(url string) (*http.Response, error) {
					if strings.Contains(url, "planets") {
						return &http.Response{
							Body: io.NopCloser(strings.NewReader(tt.searchResponse)),
						}, nil
					} else {
						return &http.Response{
							Body: io.NopCloser(strings.NewReader(tt.peopleResponse)),
						}, nil
					}
				},
			}

			output := SearchValue(tt.term, mockSwapi)

			// Check the output
			if output != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, output)
			}
		})
	}
}
