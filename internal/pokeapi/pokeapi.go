package pokeapi

import (
	"fmt"
	"io"
	"net/http"
)

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationAreaResponse struct {
	Next     *string         `json:"next"`
	Previous *string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type ExploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}
	} `json:"pokemon_encounters"`
}

type Config struct {
	Next	 *string
	Previous *string
}

var Cfg = Config{}

// Creates a new Config with the initial URL set to the first page of location areas.
func NewConfig() {
    initialURL := "https://pokeapi.co/api/v2/location-area/"
    Cfg.Next     = &initialURL
    Cfg.Previous = nil
}

// Updates the Config with the next and previous URLs from the API response.
func (cfg *Config) Update(next, previous *string) error {
    if cfg == nil {
        return fmt.Errorf("Failed to update Config; receiver is nil")
    }
    
    cfg.Next = next
    cfg.Previous = previous
    return nil
}

const (
    FORWARD = iota
    BACK
)

// Location areas are sections of areas, such as floors in a building or cave. Each area has its own set of possible Pokémon encounters.
func GetLocationAreas(direction int) ([]byte, error) {
	var resp *http.Response
	var err error

	switch direction {
	case BACK:
		if Cfg.Previous == nil {
			return []byte{}, fmt.Errorf("you're on the first page")
		}
		resp, err = http.Get(*Cfg.Previous)
	case FORWARD:
		if Cfg.Next == nil {
			return []byte{}, fmt.Errorf("Congrats! you've reached the last page. Go touch some grass now.")
		}
    	resp, err = http.Get(*Cfg.Next)
	default:
		return []byte{}, fmt.Errorf("you didn't think this through, did you?")
	}

    if err != nil {
        return []byte{}, fmt.Errorf("failed to get location areas: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return []byte{}, fmt.Errorf("failed to get location areas: status code %d", resp.StatusCode)
    }

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// var response locationAreaResponse
	// if err := json.Unmarshal(data, &response); err != nil {
	// 	return []LocationArea{}, fmt.Errorf("failed to unmarshal location areas response: %w", err)
	// }

	// pokeapi.Cfg.Update(response.Next, response.Previous)

    return data, nil
}

// GetPokemonEncounters retrieves the list of Pokémon encounters for a given location area name.
func GetPokemonEncounters(locationAreaName string) ([]byte, error) {
	locationAreaURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", locationAreaName)
	resp, err := http.Get(locationAreaURL)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to get pokemon encounters: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("failed to get pokemon encounters: status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}