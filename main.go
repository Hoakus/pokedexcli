package main

import (
	"bufio"
	"fmt"
	"github.com/Hoakus/pokedexcli/internal/api"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

type config struct {
	Previous *string
	Next     *string
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	commandStrings := ""
	for _, cmd := range commands {
		commandStrings += cmd.name + ": " + cmd.description + "\n"
	}

	fmt.Println(commandStrings)
	return nil
}

func commandMap(c *config) error {
	urlToFetch := "https://pokeapi.co/api/v2/location-area/"
	if c.Next == nil {
		if c.Previous != nil {
			fmt.Println("you're on the last page")
			return nil
		}

		c.Next = &urlToFetch
	}

	locationResponse, err := api.GetLocationArea(*c.Next)
	if err != nil {
		return err
	}

	c.Next = locationResponse.Next
	c.Previous = locationResponse.Previous
	names := locationResponse.GetResults()

	for i := range names {
		fmt.Println(names[i])
	}

	return nil
}

func commandMapB(c *config) error {
	if c.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locationResponse, err := api.GetLocationArea(*c.Previous)
	if err != nil {
		return err
	}

	c.Next = locationResponse.Next
	c.Previous = locationResponse.Previous
	names := locationResponse.GetResults()

	for i := range names {
		fmt.Println(names[i])
	}

	return nil
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 map locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the last 20 map locations",
			callback:    commandMapB,
		},
	}
}

func main() {
	c := config{}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()

			if input == "" {
				continue
			}
			cleanedInput := cleanInput(input)[0]

			if cmd, ok := commands[cleanedInput]; !ok {
				fmt.Println("Unknown command")
				continue
			} else {
				err := cmd.callback(&c)
				if err != nil {
					fmt.Printf("Error calling callback of: %v", cmd.name)
				}
			}

		}
	}
}

func cleanInput(input string) []string {
	cleanWords := []string{}
	words := strings.Split(input, " ")

	for i := range words {
		word := strings.TrimSpace(words[i])
		word = strings.ToLower(word)

		if word != "" {
			cleanWords = append(cleanWords, word)
		}
	}
	return cleanWords
}
