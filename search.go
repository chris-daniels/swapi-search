package main

import (
	"fmt"
	"strings"
)

type Entity interface {
	GetName() string
	GetPeople() []string
}

func SearchValue(term string, swapiApi SwapiApi) string {
	var output strings.Builder

	// Search planets
	fmt.Println("Searching planets...")
	planets, err := CachedSearchEntities[Planet]("planets", term, swapiApi)
	if err != nil {
		fmt.Println("Error searching planets")
		return ""
	}
	if len(planets) > 0 {
		output.WriteString(fmt.Sprintf("Found %d planet(s):\n", len(planets)))
	}
	for _, planet := range planets {
		output.WriteString(writeEntityPeople(planet, swapiApi))
	}

	// Search films
	fmt.Println("Searching films...")
	films, err := CachedSearchEntities[Film]("films", term, swapiApi)
	if err != nil {
		fmt.Println("Error searching films")
		return ""
	}
	if len(films) > 0 {
		output.WriteString(fmt.Sprintf("Found %d film(s):\n", len(films)))
	}
	for _, film := range films {
		output.WriteString(writeEntityPeople(film, swapiApi))
	}

	// Search species
	fmt.Println("Searching species...")
	species, err := CachedSearchEntities[Species]("species", term, swapiApi)
	if err != nil {
		fmt.Println("Error searching species")
		return ""
	}
	if len(species) > 0 {
		output.WriteString(fmt.Sprintf("Found %d species:\n", len(species)))
	}
	for _, specie := range species {
		output.WriteString(writeEntityPeople(specie, swapiApi))
	}

	// Search vehicles
	fmt.Println("Searching vehicles...")
	vehicles, err := CachedSearchEntities[Vehicle]("vehicles", term, swapiApi)
	if err != nil {
		fmt.Println("Error searching vehicles")
		return ""
	}
	if len(vehicles) > 0 {
		output.WriteString(fmt.Sprintf("Found %d vehicle(s):\n", len(vehicles)))
	}
	for _, vehicle := range vehicles {
		output.WriteString(writeEntityPeople(vehicle, swapiApi))
	}

	// Search starships
	fmt.Println("Searching starships...")
	starships, err := CachedSearchEntities[Starship]("starships", term, swapiApi)
	if err != nil {
		fmt.Println("Error searching starships")
		return ""
	}
	if len(starships) > 0 {
		output.WriteString(fmt.Sprintf("Found %d starship(s):\n", len(starships)))
	}
	for _, starship := range starships {
		output.WriteString(writeEntityPeople(starship, swapiApi))
	}

	// Search people
	fmt.Println("Searching people...")
	people, err := CachedSearchEntities[Person]("people", term, swapiApi)
	if err != nil {
		fmt.Println("Error searching people")
		return ""
	}
	if len(people) > 0 {
		output.WriteString(fmt.Sprintf("Found %d character(s):\n", len(people)))
	}
	for _, person := range people {
		output.WriteString(fmt.Sprintf("\t - %s\n", person.Name))
	}

	return output.String()
}

func writeEntityPeople(entity Entity, swapiApi SwapiApi) string {
	var output strings.Builder

	// Write entity name
	output.WriteString(fmt.Sprintf("\t - %s\n", entity.GetName()))

	// Write people
	for _, personUrl := range entity.GetPeople() {
		person, err := CachedFetchPerson(personUrl, swapiApi)
		if err != nil {
			return "Error fetching person"
		}

		output.WriteString(fmt.Sprintf("\t\t - %s\n", person.Name))
	}

	return output.String()
}
