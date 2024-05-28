package main

import (
	"testing"
)

// Test SearchValue function
func TestSearchValue(t *testing.T) {
	tests := []struct {
		name            string
		searchTermCache map[string]map[string]struct{}
		peopleCache     map[string]*Person
		term            string
		expected        string
	}{
		{
			name: "Verify a correctly populated cache will return a character",
			searchTermCache: map[string]map[string]struct{}{
				"luke": {"1": struct{}{}},
			},
			peopleCache: map[string]*Person{
				"1": {
					Name: "Luke Skywalker",
				},
			},
			term:     "luke",
			expected: "Found 1 character(s):\n\tLuke Skywalker\n",
		},
		{
			name: "Verify no characters are found when bad search term is given",
			searchTermCache: map[string]map[string]struct{}{
				"luke": {"1": struct{}{}},
			},
			peopleCache: map[string]*Person{
				"1": {
					Name: "Luke Skywalker",
				},
			},
			term:     "this character does not exist",
			expected: "No results found\n",
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function
			result := SearchValue(tt.term, tt.searchTermCache, tt.peopleCache)

			// Check the result
			if result != tt.expected {
				t.Errorf("SearchValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Lazy integration test for TestFetchSwapiData
func TestFetchSwapiData(t *testing.T) {
	searchTermCache := make(map[string]map[string]struct{})
	peopleCache := make(map[string]*Person)

	err := FetchSwapiData(peopleCache, searchTermCache)
	if err != nil {
		t.Errorf("FetchSwapiData() = %v, want nil", err)
	}

	// Assert that the cache is populated
	if len(peopleCache) != 82 {
		t.Errorf("len(peopleCache) = %v, want 82", len(peopleCache))
	}
	if len(searchTermCache) != 178 {
		t.Errorf("len(searchTermCache) = %v, want 178", len(searchTermCache))
	}
}
