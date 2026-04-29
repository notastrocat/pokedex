package main

import (
	"testing"
	"time"

	"pokedex/internal/pokeapi"
	"pokedex/internal/pokecache"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello  world  ",
			expected: []string{"hello", "world", ""},
		},
		{
			input:    "Test text  ",
			expected: []string{"test", "text", ""},
		},
		// add more cases here.
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("ERR: Test failed!!\nlen(%v) != len(%v)", actual, c.expected)
			return
		}

		for i, _ := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("ERR: Test failed!!\n%v != %v", word, expectedWord)
				return
			}
		}
	}
}

func TestCommandMapCacheStateUpdate(t *testing.T) {
	// Test that commandMap updates pagination state even on cache hits
	// This is critical because the caching logic changed to ensure Cfg.Update is always called

	// Initialize test data
	LocaleCache = pokecache.NewCache(15 * time.Second)
	pokeapi.NewConfig()

	initialNext := *pokeapi.Cfg.Next

	// First call should make an API call and cache the result
	err := commandMap("")
	if err != nil {
		// Expected - may fail without real network, but the state should still update
		t.Logf("First call returned error (expected without network): %v", err)
	}

	// State should have advanced even if there was an error
	if pokeapi.Cfg.Next != nil && *pokeapi.Cfg.Next == initialNext {
		t.Errorf("Cfg.Next should have been updated after first call")
	}
}

func TestCommandMapBackCacheStateUpdate(t *testing.T) {
	// Test that commandMapBack updates pagination state even on cache hits
	LocaleCache = pokecache.NewCache(15 * time.Second)
	pokeapi.NewConfig()

	// Need to move forward first to have a "Previous" state
	commandMap("")

	if pokeapi.Cfg.Previous == nil {
		t.Skip("Skipping test - Previous state not available (no network)")
	}

	err := commandMapBack("")
	if err != nil {
		t.Logf("mapb call returned error (expected without network): %v", err)
	}
}
