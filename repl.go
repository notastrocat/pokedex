package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return fmt.Errorf("cannot exit Pokedex.")
}

var CommandRegistry map[string]cliCommand

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")

	for _, cmd := range CommandRegistry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
    fmt.Println()

	return nil
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cmd := cleanInput(input)
		if cmd[0] == "" {
			continue
		}
		if command, exists := CommandRegistry[cmd[0]]; exists {
			err := command.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command:", cmd)
		}
	}
}
