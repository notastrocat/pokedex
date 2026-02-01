package main

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func init() {
	CommandRegistry = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
        "map" : {
            name:        "map",
            description: "Displays a list of locations, 20 at a time.",
            callback:    commandMap,
        },
        "mapb": {
            name:        "mapb",
            description: "Displays the previous 20 locations.",
            callback:    commandMapBack,
        },
	}
}

func main() {
    repl()
}
