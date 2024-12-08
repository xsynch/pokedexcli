package main

import (
	"bufio"
	"strings"
	"time"

	"fmt"
	"log"
	"os"
)

type Command interface {
	Execute(cfg *ConfigPokeAreas, args ...string) error
}

type commandHelpStruct struct{}

func (h commandHelpStruct) Execute(cfg *ConfigPokeAreas, args ...string) error {
	fmt.Printf("Welcome to Pokedex\nUsage:\n\nhelp: Display a help message(this)\nexit: Exit the Pokedex\n\n")
	return nil
}

type commandExitStruct struct{}

func (e commandExitStruct) Execute(cfg *ConfigPokeAreas, args ...string) error {
	os.Exit(0)
	return nil
}

type commandMapstruct struct{}

func (m commandMapstruct) Execute(cfg *ConfigPokeAreas, args ...string) error {
	if cfg.Next == nil {
		cfg.Next = &APIENDPOINT
	}

	err := cfg.getPokeLocations(*cfg.Next)
	if err != nil {

		return fmt.Errorf("error getting areas information: %s", err)
	}
	return nil

}

type commaandMapbstruct struct{}

func (mb commaandMapbstruct) Execute(cfg *ConfigPokeAreas, args ...string) error {
	err := cfg.getPreviousPokeAreas()
	if err != nil {
		return fmt.Errorf("error getting previous areas: %s", err)
	}
	return nil

}

type commaandExplorestruct struct{}

func (ex commaandExplorestruct) Execute(cfg *ConfigPokeAreas, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("explore takes two arguments: explore <area>")
	}
	log.Printf("Exploring %s...", args[0])
	err := cfg.exploreArea(args[0])
	if err != nil {
		return err
	}
	return nil
}

type commandCatch struct{}

func (catch commandCatch) Execute(cfg *ConfigPokeAreas, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("catch takes two arguments: catch <pokemon_name>")
	}
	err := cfg.catchPokemon(args[0])
	if err != nil {
		return err
	}
	return nil
}

type commandInspect struct{}

func (catch commandInspect) Execute(cfg *ConfigPokeAreas, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("catch takes two arguments: catch <pokemon_name>")
	}
	err := cfg.inspectPokemon(args[0])
	if err != nil {
		return err
	}
	return nil
}


type commandPokedex struct{}

func (pok commandPokedex) Execute(cfg *ConfigPokeAreas, args ...string) error {

	err := cfg.printPokedex()
	if err != nil {
		return err
	}
	return nil
}

func commandHelp(cfg *ConfigPokeAreas) error {

	fmt.Printf("Welcome to Pokedex\nUsage:\n\nhelp: Display a help message(this)\nexit: Exit the Pokedex\n\n")
	return nil
}

func exitCommand(cfg *ConfigPokeAreas) error {
	os.Exit(0)
	return nil
}

func prompt() {
	fmt.Print("pokedex > ")
}

type mainMenu struct {
	name        string
	description string
	// function func(*ConfigPokeAreas)error
	command Command
}

func NewMainMenu() map[string]mainMenu {
	nm := map[string]mainMenu{
		"help": {
			name:        "help",
			description: "Displays help message",
			command:     commandHelpStruct{},
		},
		"exit": {
			name:        "exit",
			description: "Exits the program",
			command:     commandExitStruct{},
		},
		"map": {
			name:        "map",
			description: "Show 20 Pokemon area names",
			command:     commandMapstruct{},
		},
		"mapb": {
			name:        "map",
			description: "Show previous 20 Pokemon area names",
			command:     commaandMapbstruct{},
		},
		"explore": {
			name:        "explore",
			description: "Explore Pokemon in the area",
			command:     commaandExplorestruct{},
		},
		"catch": {
			name:        "catch",
			description: "Try to catch the pokemon requested",
			command:     commandCatch{},
		},
		"inspect": {
			name:        "catch",
			description: "Get the stats of a caught pokemon",
			command:     commandInspect{},
		},
		"pokedex": {
			name:        "pokedex",
			description: "Print your current pokedex",
			command:     commandPokedex{},
		},
	}
	return nm
}

type res struct {
	Name string `json:"name"`
	URL  string `json:"url"`
} //`json:"results"`

func main() {

	prompt()
	var line string

	newMenu := NewMainMenu()
	expiration := time.Second * 25
	cfg := NewPokeCli(expiration)
	// cfg2 := &ConfigPokeAreas{
	// 	Count: 0,
	// 	Next: nil,
	// 	Previous: nil,

	// }

	for {

		scanner := bufio.NewScanner(os.Stdin)
		var args []string

		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			log.Printf("%s\n", err)
			// os.Exit(1)
		}
		line = scanner.Text()
		if len(line) > 1 {
			parts := strings.Fields(line)
			// command := parts[0]
			args = parts[1:]
			line = parts[0]

		}

		if value, ok := newMenu[line]; ok {
			switch line {
			case "explore":
				err := value.command.Execute(cfg, args...)
				if err != nil {

					log.Printf("%s\n", err)

				}
			case "catch":
				err := value.command.Execute(cfg, args...)
				if err != nil {
					log.Printf("%s\n", err)
				}
			case "inspect":
				err := value.command.Execute(cfg, args...)
				if err != nil {
					log.Printf("%s\n", err)
				}

			default:
				err := value.command.Execute(cfg, "")
				if err != nil {
					log.Printf("%s\n", err)

				}

			}
			// err := value.function(cfg)
			// if err != nil {
			// 	log.Printf("%s\n",err)

			// }
			prompt()

		} else {
			_ = commandHelp(cfg)
			prompt()

		}

	}
}
