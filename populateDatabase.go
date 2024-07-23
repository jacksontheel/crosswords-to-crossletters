package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type Crossword struct {
	Answers AnswerClues
	Clues   AnswerClues
	Author  string
}

type AnswerClues struct {
	Across []string
	Down   []string
}

func PopulateDatabase(database Database, crosswordDirectory, commonWordsFile string) {
	var wg sync.WaitGroup

	// err := addClues(wg, database, "./data/crosswords")
	err := addClues(wg, database, crosswordDirectory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// err = addWords(wg, database, "./data/common_words.txt")
	err = addWords(wg, database, commonWordsFile)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

func addClues(wg sync.WaitGroup, db Database, directory string) error {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(info.Name(), "json") {
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return err
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return err
			}

			var crossword Crossword

			err = json.Unmarshal([]byte(data), &crossword)
			if err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				return err
			}

			for i := range crossword.Clues.Across {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					insertClueIntoDatabase(db, crossword.Clues.Across[i], crossword.Answers.Across[i], crossword.Author)
				}(i)

			}

			for i := range crossword.Clues.Down {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					insertClueIntoDatabase(db, crossword.Clues.Down[i], crossword.Answers.Down[i], crossword.Author)
				}(i)
			}

			wg.Wait()
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func addWords(wg sync.WaitGroup, db Database, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		err := db.addWord(strings.ToUpper(scanner.Text()))
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	wg.Wait()

	return nil
}

func insertClueIntoDatabase(db Database, clue, answer, author string) {
	// Remove beginning number and decimal from clue, leaving only the clue text
	clueRE := regexp.MustCompile(`^\d+\.\s*`)
	clueText := clueRE.ReplaceAllString(clue, "")

	// Only add clue if clue text doesn't contain something like "10-across"
	downAcrossRE := regexp.MustCompile(`\d+-(Across|across|Down|down)`)

	if !downAcrossRE.Match([]byte(clueText)) {
		id, _ := db.addClue(clueText, answer, author)

		for _, r := range answer {
			db.addLetter(id, r)
		}
	} else {
		fmt.Printf("Excluding %s\n", clueText)
	}
}
