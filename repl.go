package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"pokedex/internal/pokeapi"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	words = append(words, "") // add an empty string to the end of the slice to prevent out-of-range errors when accessing cmd[1]
	return words
}

func commandExit(locationAreaName string) error {
	log.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return fmt.Errorf("cannot exit Pokedex.")
}

var CommandRegistry map[string]cliCommand

func commandHelp(locationAreaName string) error {
	log.Println("Welcome to the Pokedex!")
	log.Println("Usage:")
	log.Println()

	for _, cmd := range CommandRegistry {
		log.Printf("%7s: %s\n", cmd.name, cmd.description)
	}
	log.Println()

	return nil
}

func commandMap(locationAreaName string) error {
	// if pokeapi.Cfg.Previous != nil && pokeapi.Cfg.Next != nil {
	// 	log.Println(*pokeapi.Cfg.Next)
	// 	log.Println(*pokeapi.Cfg.Previous)
	// }
	if pokeapi.Cfg.Next == nil {
		return fmt.Errorf("don't know about next URL.")
	}

	var data []byte
	cachedData, ok := LocaleCache.Get(*pokeapi.Cfg.Next)
	if ok {
		// cache hit! we can skip the API call
		data = cachedData
	} else {
		// call the func that handles the API call to get the Location Areas...
		var err error
		data, err = pokeapi.GetLocationAreas(pokeapi.FORWARD)
		if err != nil {
			return err
		}

		if pokeapi.Cfg.Next != nil {
			log.Println("updating the cache for key", *pokeapi.Cfg.Next)
			LocaleCache.Add(*pokeapi.Cfg.Next, data)
		}
	}

	var response pokeapi.LocationAreaResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("failed to unmarshal location areas response: %w", err)
	}

	pokeapi.Cfg.Update(response.Next, response.Previous)

	for _, area := range response.Results {
		log.Printf("%s\n", area.Name)
	}

	// fmt.Println("--------------------------------------------------------------------")

	// log.Println(LocaleCache)

	// if pokeapi.Cfg.Previous != nil && pokeapi.Cfg.Next != nil {
	// 	log.Println(*pokeapi.Cfg.Next)
	// 	log.Println(*pokeapi.Cfg.Previous)
	// }

	return nil
}

func commandMapBack(locationAreaName string) error {
	// if pokeapi.Cfg.Previous != nil && pokeapi.Cfg.Next != nil {
	// 	log.Println(*pokeapi.Cfg.Next)
	// 	log.Println(*pokeapi.Cfg.Previous)
	// }
	if pokeapi.Cfg.Previous == nil {
		return fmt.Errorf("you're on the first page")
	}

	var data []byte
	cachedData, ok := LocaleCache.Get(*pokeapi.Cfg.Previous)
	if ok {
		// cache hit! we can skip the API call
		data = cachedData
	} else {
		// call the func that handles the API call to get the previous Location Areas...
		var err error
		data, err = pokeapi.GetLocationAreas(pokeapi.BACK)
		if err != nil {
			return err
		}

		if pokeapi.Cfg.Previous != nil {
			log.Println("updating the cache for key", *pokeapi.Cfg.Previous)
			LocaleCache.Add(*pokeapi.Cfg.Previous, data)
		}
	}

	var response pokeapi.LocationAreaResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("failed to unmarshal location areas response: %w", err)
	}

	pokeapi.Cfg.Update(response.Next, response.Previous)

	// log.Println("Previous Location Areas:")
	for _, area := range response.Results {
		log.Printf("%s\n", area.Name)
	}

	// if pokeapi.Cfg.Previous != nil && pokeapi.Cfg.Next != nil {
	// 	log.Println(*pokeapi.Cfg.Next)
	// 	log.Println(*pokeapi.Cfg.Previous)
	// }

	return nil
}

func commandExplore(locationAreaName string) error {
	log.Printf("Exploring %s...\n", locationAreaName)

	cachedData, ok := PokemonCache.Get(locationAreaName)
	if ok {
		// cache hit! we can skip the API call and just print the cached pokemon encounters.
		var cachedEncounters pokeapi.ExploreResponse
		err := json.Unmarshal(cachedData, &cachedEncounters)
		if err != nil {
			return fmt.Errorf("failed to unmarshal cached data: %v", err)
		}

		log.Println("Found Pokemon:")
		for _, encounter := range cachedEncounters.PokemonEncounters {
			log.Printf(" - %s\n", encounter.Pokemon.Name)
		}

		return nil
	}

	data, err := pokeapi.GetPokemonEncounters(locationAreaName)
	if err != nil {
		return err
	}

	var response pokeapi.ExploreResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("failed to unmarshal pokemon encounters response: %w", err)
	}

	log.Println("Found Pokemon:")
	for _, encounter := range response.PokemonEncounters {
		log.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	PokemonCache.Add(locationAreaName, data)

	return nil
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)

	pokeapi.NewConfig()

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cmd := cleanInput(input)
		if cmd[0] == "" {
			continue
		}
		if command, exists := CommandRegistry[cmd[0]]; exists {
			err := command.callback(cmd[1])
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println("Unknown command:", cmd)
		}
	}
}
