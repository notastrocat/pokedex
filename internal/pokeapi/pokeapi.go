package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type locationAreaResponse struct {
	Next     *string         `json:"next"`
	Previous *string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type Config struct {
	Next	 *string
	Previous *string
}

var Cfg = Config{}

func NewConfig() {
    initialURL := "https://pokeapi.co/api/v2/location-area/"
    Cfg.Next     = &initialURL
    Cfg.Previous = nil
}

func (cfg *Config) update(next, previous *string) error {
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
func GetLocationAreas(direction int) ([]LocationArea, error) {
	var resp *http.Response
	var err error

	switch direction {
	case BACK:
		if Cfg.Previous == nil {
			return []LocationArea{}, fmt.Errorf("you're on the first page")
		}
		resp, err = http.Get(*Cfg.Previous)
	case FORWARD:
		if Cfg.Next == nil {
			return []LocationArea{}, fmt.Errorf("Congrats! you've reached the last page. Go touch some grass now.")
		}
    	resp, err = http.Get(*Cfg.Next)
	default:
		return []LocationArea{}, fmt.Errorf("you didn't think this through, did you?")
	}

    if err != nil {
        return []LocationArea{}, fmt.Errorf("failed to get location areas: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return []LocationArea{}, fmt.Errorf("failed to get location areas: status code %d", resp.StatusCode)
    }

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []LocationArea{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var response locationAreaResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return []LocationArea{}, fmt.Errorf("failed to unmarshal location areas response: %w", err)
	}

	Cfg.update(response.Next, response.Previous)

    return response.Results, nil
}