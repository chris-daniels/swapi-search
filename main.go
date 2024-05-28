package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var swapi SwapiApi = SwapiApi{
	handler: func(url string) (*http.Response, error) {
		return http.Get(url)
	},
}

func main() {
	// Preload people cache
	fmt.Println("Preloading people cache...")
	err := PreFetchPeople()
	if err != nil {
		fmt.Println("Error preloading people cache")
	}

	// Loop to accept search terms
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter search term: ")
		term, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading search term")
			return
		}

		term = strings.TrimSpace(term)
		if term == "" {
			break
		}

		output := SearchValue(term, swapi)

		fmt.Println("*********************************")
		fmt.Println("Search results")
		fmt.Println("*********************************")
		fmt.Println(output)
	}
}
