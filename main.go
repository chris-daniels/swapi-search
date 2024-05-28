package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Entity interface {
	GetName() string
	GetPeople() []string
}

func populateSearchTerms(entity Entity, searchTermCache map[string]map[string]struct{}) {
	term := strings.ToLower(entity.GetName())
	personIds, ok := searchTermCache[term]
	if !ok {
		personIds = map[string]struct{}{}
	}

	for _, personUrl := range entity.GetPeople() {
		personId := getIdFromUrl(personUrl)
		personIds[personId] = struct{}{}
	}

	searchTermCache[term] = personIds
}

func FetchSwapiData(peopleCache map[string]*Person, searchTermCache map[string]map[string]struct{}) error {
	// Fetch planets and populate search terms
	planets, err := FetchEntities[Planet]("planets")
	if err != nil {
		fmt.Println("Error fetching planets: ", err)
		return err
	}
	for _, planet := range planets {
		populateSearchTerms(planet, searchTermCache)
	}

	// Fetch films and populate search terms
	films, err := FetchEntities[Film]("films")
	if err != nil {
		fmt.Println("Error fetching films: ", err)
		return err
	}
	for _, film := range films {
		populateSearchTerms(film, searchTermCache)
	}

	// Fetch species and populate search terms
	species, err := FetchEntities[Species]("species")
	if err != nil {
		fmt.Println("Error fetching species: ", err)
		return err
	}
	for _, specie := range species {
		populateSearchTerms(specie, searchTermCache)
	}

	// Fetch vehicles and populate search terms
	vehicles, err := FetchEntities[Vehicle]("vehicles")
	if err != nil {
		fmt.Println("Error fetching vehicles: ", err)
		return err
	}
	for _, vehicle := range vehicles {
		populateSearchTerms(vehicle, searchTermCache)
	}

	// Fetch starships and populate search terms
	starships, err := FetchEntities[Starship]("starships")
	if err != nil {
		fmt.Println("Error fetching starships: ", err)
		return err
	}
	for _, starship := range starships {
		populateSearchTerms(starship, searchTermCache)
	}

	// Fetch people and populate people cache
	people, err := FetchEntities[Person]("people")
	if err != nil {
		fmt.Println("Error fetching people: ", err)
		return err
	}
	for _, person := range people {
		peopleCache[getIdFromUrl(person.URL)] = person
	}

	return nil
}

func SearchValue(term string, searchTermCache map[string]map[string]struct{}, peopleCache map[string]*Person) string {
	personIds, ok := searchTermCache[strings.ToLower(term)]
	if !ok {
		return "No results found\n"
	}

	output := fmt.Sprintf("Found %d character(s):\n", len(personIds))

	for personId := range personIds {
		person, ok := peopleCache[personId]
		if !ok {
			fmt.Println("Person not found")
			continue
		}

		output += fmt.Sprintf("\t%s\n", person.Name)
	}
	return output
}

func main() {
	// Store mapping of id to Person
	var peopleCache = map[string]*Person{}

	// Store map of search term to set of people ids
	var searchTermCache = map[string]map[string]struct{}{}

	// Pre-populate search term cache
	fmt.Println("Fetching swapi data...")
	err := FetchSwapiData(peopleCache, searchTermCache)
	if err != nil {
		fmt.Println("Error fetching swapi data")
		return
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

		output := SearchValue(term, searchTermCache, peopleCache)
		fmt.Println(output)
	}
}
