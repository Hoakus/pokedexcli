package main

import (
	"bufio"
	"fmt"
	"github.com/Hoakus/pokedexcli/internal/api"
	"math/rand"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, arg string) error
}

type config struct {
	Previous *string
	Next     *string
}

const (
	locationAreaURL = "https://pokeapi.co/api/v2/location-area/"
	pokemonURL      = "https://pokeapi.co/api/v2/pokemon/"
)

var caughtPokemon map[string]api.PokemonResponse

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
		"explore": {
			name:        "explore",
			description: `Displays a list of all pokemon in the area | explore 'area-name'`,
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: `Try to catch a pokemon! | catch 'pokemon-name'`,
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: `Inspect a caught pokemon | inspect 'pokemon-name'`,
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: `View all pokemon in your pokedex`,
			callback:    commandPokedex,
		},
	}
}

// Pokedex Functions
func commandExit(c *config, _ string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, _ string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	commandStrings := ""
	for _, cmd := range commands {
		commandStrings += cmd.name + ": " + cmd.description + "\n"
	}

	fmt.Println(commandStrings)
	return nil
}

func commandMap(c *config, _ string) error {
	baseUrl := locationAreaURL
	if c.Next == nil {
		if c.Previous != nil {
			fmt.Println("you're on the last page")
			return nil
		} else {
			c.Next = &baseUrl
		}
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

func commandMapB(c *config, _ string) error {
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

	areaNames := locationResponse.GetResults()
	for i := range areaNames {
		fmt.Println(areaNames[i])
	}
	return nil
}

func commandExplore(c *config, areaName string) error {
	if areaName == "" {
		fmt.Printf("enter an area name | explore 'area name'")
		return nil
	}

	fullUrl := locationAreaURL + areaName + "/"
	fmt.Printf("Exploring %v...\n", areaName)

	exploreResponse, err := api.GetAreaByName(fullUrl)
	if err != nil {
		fmt.Printf("could not explore %v\n", areaName)
		return err
	}

	fmt.Println("Found Pokemon:")
	pokeNames := exploreResponse.GetResults()
	for i := range pokeNames {
		fmt.Printf("- %v\n", pokeNames[i])
	}
	return nil
}

func commandCatch(c *config, pokemonName string) error {
	if pokemonName == "" {
		fmt.Println("enter a pokemon name | catch 'pokemon name'")
		return nil
	}

	fullUrl := pokemonURL + pokemonName + "/"

	pokemonResponse, err := api.GetPokemonByName(fullUrl)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)
	catchAttempt := rand.Intn(100)
	if catchAttempt < (pokemonResponse.BaseXp / 5) {
		fmt.Printf("%v escaped!\n", pokemonName)
		return nil
	}

	fmt.Printf("%v was caught!\n", pokemonName)
	fmt.Println("you can now inspect it with 'inspect'")

	if caughtPokemon == nil {
		caughtPokemon = make(map[string]api.PokemonResponse)
	}

	caughtPokemon[pokemonName] = pokemonResponse

	return nil
}

func commandInspect(c *config, pokemonName string) error {
	var pokemon api.PokemonResponse
	if pokemonName == "" {
		fmt.Println("enter a pokemon name | catch 'pokemon name'")
		return nil
	}

	if caughtPokemon == nil {
		fmt.Println("you haven't caught any Pokemon")
	}

	if _, ok := caughtPokemon[pokemonName]; !ok {
		fmt.Printf("%v hasn't been caught yet\n", pokemonName)
	}

	pokemon = caughtPokemon[pokemonName]

	statsOutput := ""
	for i := range pokemon.Stats {
		stat := pokemon.Stats[i]
		total := stat.BaseStat + stat.Effort
		statsOutput += fmt.Sprintf(" -%v: %d\n", stat.Stat.Name, total)
	}

	typesOutput := ""
	for i := range pokemon.Types {
		typesOutput += fmt.Sprintf(" - %v\n", pokemon.Types[i].Type.Name)
	}

	fmt.Printf("Name: %v\nHeight: %v\nWeight: %v\nStats:\n%vTypes:\n%v", pokemonName, pokemon.Height, pokemon.Weight, statsOutput, typesOutput)
	return nil
}

func commandPokedex(c *config, _ string) error {
	if caughtPokemon == nil {
		fmt.Println("you don't have any pokemon in your pokedex")
		return nil
	}

	fmt.Println("Your Pokemon:")
	for pokemon, _ := range caughtPokemon {
		fmt.Printf(" -%v\n", pokemon)
	}
	return nil
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
			cleanedArgs := cleanInput(input)

			if _, ok := commands[cleanedArgs[0]]; !ok {
				fmt.Println("Unknown command")
				continue
			}

			cmd := commands[cleanedArgs[0]]
			arg := ""
			if len(cleanedArgs) > 1 {
				arg = cleanedArgs[1]
			}
			err := cmd.callback(&c, arg)
			if err != nil {
				fmt.Printf("Error calling callback of: %v", cmd.name)

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
