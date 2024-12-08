package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"time"

	pokecache "github.com/xsynch/pokedexcli/internal"
)

var APIENDPOINT = "https://pokeapi.co/api/v2/location-area"
var APIENDPOINT_NEXT = "https://pokeapi.co/api/v2/location-area"
var APIENDPOINT_PREV = ""
var CHARACTER_API_ENDPOINT = "https://pokeapi.co/api/v2/pokemon"

type ConfigPokeAreas struct {
	Count              int     `json:"count"`
	Next               *string `json:"next"`
	Previous           *string `json:"previous"`
	Results            []res
	PokeAreaInformtion []PokeAreasDetails
	pokeCache          pokecache.PokeCache
	pokeMon            map[string]PokeMonCharacter
}

type PokeAreasDetails struct {
	Areas []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"areas"`
	GameIndices []struct {
		GameIndex  int `json:"game_index"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"game_indices"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	Region struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"region"`
}

type ExploredArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type PokeMonCharacter struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Height    int `json:"height"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func NewPokeCli(expiration time.Duration) *ConfigPokeAreas {
	return &ConfigPokeAreas{
		Count:     0,
		Next:      nil,
		Previous:  nil,
		pokeCache: pokecache.NewCache(expiration),
		pokeMon:   map[string]PokeMonCharacter{},
	}
}

func (c *ConfigPokeAreas) checkCache(url string) ([]byte, error) {
	_, ok := c.pokeCache.Entry[url]
	var body []byte
	if !ok {

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading the request: %s", err)
			return nil, err
		}
		log.Printf("Adding to the cache url: %s", url)
		c.pokeCache.Add(url, body)
		return body, nil
	} else {
		log.Printf("Cache Used, found %s", url)
		body, _ = c.pokeCache.Get(url)
	}

	return body, nil
}

func (c *ConfigPokeAreas) getPokeLocations(url string) error {

	pa := ConfigPokeAreas{}

	body, err := c.checkCache(url)
	if err != nil {
		log.Printf("Error reading the body\n")
		return fmt.Errorf("error reading the body")
	}

	err = json.Unmarshal(body, &pa)
	if err != nil {
		log.Printf("Error unmarshalling %s: %s", body, err)
		return err
	}
	// log.Printf("Count: %d and Next: %s\n\n", pa.Count, pa.Next)
	for _, val := range pa.Results {
		// fmt.Printf("Pa Name: %s with value: %s\n", val.Name, val.URL)
		// c.getPokeAreaDetails(val.URL)

		fmt.Printf("%s\n", val.Name)

	}
	c.Previous = pa.Previous
	c.Next = pa.Next
	log.Printf("next location set to %v\n", *c.Next)

	return nil

	// stringBody := string(body)
	// fmt.Println(stringBody)
}

func (c *ConfigPokeAreas) getPreviousPokeAreas() error {

	if c.Previous == nil {
		return fmt.Errorf("no previous areas to return")
	}

	log.Printf("Trying to set the api endpoint to %s", *c.Previous)
	//c.Next = c.Previous
	err := c.getPokeLocations(*c.Previous)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigPokeAreas) setCount() error {
	urlNext, err := url.Parse(APIENDPOINT_NEXT)
	if err != nil {
		log.Printf("Error parsing the URL: %s", err)
		return err
	}
	params, _ := url.ParseQuery(urlNext.RawQuery)
	if _, ok := params["offset"]; ok {
		result, err := strconv.Atoi(params["offset"][0])
		if err != nil {
			log.Printf("There was an error converting %v to an integer\n", params["offset"])
			return err
		}
		log.Printf("Change previous from %d to %d", result, result-20)
		c.Count = result - 20
	} else {
		c.Count = 0
		APIENDPOINT_PREV = APIENDPOINT
	}
	return nil

}

func (c *ConfigPokeAreas) exploreArea(location string) error {
	url := fmt.Sprintf("%s/%s", APIENDPOINT, location)
	// log.Printf("Searching %s with url: %s", location, url)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Error getting data from %s", url)
		return err
	}

	if response.StatusCode == 404 {

		return errors.New("area not found")
	}

	body, err := c.checkCache(url)
	if err != nil {
		log.Printf("Error reading the body\n")
		return fmt.Errorf("error reading the body")
	}

	area_details := ExploredArea{}
	// body, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	log.Printf("Error reading the response body\n")
	// 	return err
	// }

	err = json.Unmarshal(body, &area_details)
	if err != nil {
		return err
	}
	for _, val := range area_details.PokemonEncounters {
		fmt.Printf("- %s\n", val.Pokemon.Name)
	}

	return errors.New("testing error return")
}

func (c *ConfigPokeAreas) catchPokemon(name string) error {
	var max_experience = 500
	fmt.Printf("Throwing a ball at %s...\n", name)
	url := fmt.Sprintf("%s/%s", CHARACTER_API_ENDPOINT, name)
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error getting data from url: %s with error %s", url, err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusNotFound {
		return fmt.Errorf("pokemon character %s not found", name)
	}
	pokemon_char := PokeMonCharacter{}
	err = json.Unmarshal(body, &pokemon_char)
	if err != nil {
		return err
	}
	var probability_to_catch = int((1 - (pokemon_char.BaseExperience / max_experience)) * 100)
	r := rand.IntN(100)
	if r > probability_to_catch {
		fmt.Printf("%s has escaped.\n", name)
	} else {
		fmt.Printf("%s was caught\n", name)
		c.pokeMon[name] = pokemon_char
	}
	//  log.Printf("probability: %d random number: %d base exp: %d",probability_to_catch, r, pokemon_char.BaseExperience)

	return nil
}

func (c *ConfigPokeAreas) inspectPokemon(name string) error {
	val, ok := c.pokeMon[name]
	if ok {
		log.Printf("Found %s within the pokedex", name)
	} else {
		return fmt.Errorf("you have not caught %s yet.", name)		
		
	}
	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats: \n", name, val.Height, val.Weight)
	for _,value := range val.Stats {
		fmt.Printf("\t-%s:%d\n", value.Stat.Name,value.BaseStat)
		
	}
	fmt.Printf("Types:\n")
	for _,value := range val.Types {
		fmt.Printf("\t- %s\n",value.Type.Name)
	}

	return nil
}
