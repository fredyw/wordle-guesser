package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type constraint struct {
	correctSpots map[int]map[string]bool
	wrongSpots   map[int]map[string]bool
	invalid      map[string]bool
}

func guessWords(dictionaryPath string, constraint constraint) []string {
	file, err := os.Open(dictionaryPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var guessedWords []string
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) == 0 {
			continue
		}
		validChars := validChars(constraint.wrongSpots)
		count := 0
		for _, c := range word {
			if _, found := validChars[string(c)]; found {
				count++
			}
		}
		// Make sure that the word contains all valid characters.
		if count < len(validChars) {
			continue
		}
		possibleWord := true
		for i, c := range word {
			if _, found := constraint.invalid[string(c)]; found {
				possibleWord = false
				break
			}
			if correctSpot, found := constraint.correctSpots[i+1]; found {
				if _, found := correctSpot[string(c)]; !found {
					possibleWord = false
					break
				}
			}
			if wrongSpot, found := constraint.wrongSpots[i+1]; found {
				if _, found := wrongSpot[string(c)]; found {
					possibleWord = false
					break
				}
			}
		}
		if possibleWord {
			guessedWords = append(guessedWords, word)
		}
	}
	return guessedWords
}

func validChars(wrongSpots map[int]map[string]bool) map[string]bool {
	validChars := map[string]bool{}
	for _, m := range wrongSpots {
		for c := range m {
			validChars[c] = true
		}
	}
	return validChars
}

var dictionaryFlag string
var correctSpotFlag string
var wrongSpotFlag string
var invalidFlag string

func init() {
	flag.StringVar(
		&dictionaryFlag,
		"dictionary",
		"sgb-words.txt",
		"Path to the dictionary file.")
	flag.StringVar(
		&correctSpotFlag,
		"correct-spot",
		"",
		"Characters in the correct spots.\n"+
			"Format : <position1>:<characters>;<position2>:<characters>,...\n"+
			"Example: 1:e;2:p;3:o")
	flag.StringVar(
		&wrongSpotFlag,
		"wrong-spot",
		"",
		"Characters in the wrong spots.\n"+
			"Format : <position1>:<characters>;<position2>:<characters>,...\n"+
			"Example: 2:e;3:p,e;4:o")
	flag.StringVar(
		&invalidFlag,
		"invalid",
		"",
		"Invalid characters.\n"+
			"Format : <chars>\n"+
			"Example: t,a,s,d")
}

func validateFlags() {
	if _, err := os.Stat(dictionaryFlag); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Error:", dictionaryFlag, "does not exist")
		os.Exit(1)
	}
	if len(correctSpotFlag) > 0 {
		validateCharSpotFlag(correctSpotFlag)
	}
	if len(wrongSpotFlag) > 0 {
		validateCharSpotFlag(wrongSpotFlag)
	}
}

func validateCharSpotFlag(charSpotFlag string) {
	for _, p := range strings.Split(charSpotFlag, ";") {
		split := strings.Split(p, ":")
		if len(split) != 2 {
			flag.Usage()
			os.Exit(1)
		}
	}
}

func buildConstraint() (constraint, error) {
	correctSpots, err := buildCharSpotConstraint(correctSpotFlag)
	if err != nil {
		return constraint{}, err
	}
	wrongSpots, err := buildCharSpotConstraint(wrongSpotFlag)
	if err != nil {
		return constraint{}, err
	}
	invalid := map[string]bool{}
	for _, c := range strings.Split(invalidFlag, ",") {
		invalid[c] = true
	}
	return constraint{
		correctSpots: correctSpots,
		wrongSpots:   wrongSpots,
		invalid:      invalid,
	}, nil
}

func buildCharSpotConstraint(charSpotFlag string) (map[int]map[string]bool, error) {
	charSpots := map[int]map[string]bool{}
	if len(charSpotFlag) > 0 {
		for _, p := range strings.Split(charSpotFlag, ";") {
			split := strings.Split(p, ":")
			position := split[0]
			chars := strings.Split(split[1], ",")
			m := map[string]bool{}
			for _, c := range chars {
				m[c] = true
			}
			p, err := strconv.Atoi(position)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Invalid position %s", position))
			}
			charSpots[p] = m
		}
	}
	return charSpots, nil
}

func main() {
	flag.Parse()
	validateFlags()
	constraint, err := buildConstraint()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	words := guessWords(dictionaryFlag, constraint)
	fmt.Println("Possible words:")
	for _, word := range words {
		fmt.Println("-", word)
	}
}
