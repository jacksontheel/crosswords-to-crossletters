package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func CreatePuzzles(database Database, numQuestions int, outputDirectory string) {
	puzzlesCreated := 0
	usedQuestions := make(map[Question]interface{}, 0)

	for puzzlesCreated < 365 {
		letters := append(
			TakeRandomSubset(Vowels, 3),
			TakeRandomSubset(Consonants, 4)...,
		)

		allQuestions, err := database.queryWordsWithLetters(letters[0], letters[1], letters[2], letters[3], letters[4], letters[5], letters[6])

		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		result := Puzzle{
			Letters:   letters,
			Questions: allQuestions,
		}

		result = result.
			FilterQuestions(func(q Question) bool {
				_, ok := usedQuestions[q]
				return !ok
			}).
			FilterDuplicateQuestionAnswers().
			WithRandomSubsetOfQuestions(numQuestions)

		if len(result.Questions) < numQuestions {
			continue
		}

		// Add the questions of the resulting puzzle to the set, so they are not reused
		for _, q := range result.Questions {
			usedQuestions[q] = struct{}{}
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			// todo
			log.Fatalf("Error marshaling struct: %v", err)
		}

		file, err := os.Create(fmt.Sprintf("%s/puzzle%d.json", outputDirectory, puzzlesCreated))
		if err != nil {
			fmt.Println("Error creating file:", err)
			continue
		}
		defer file.Close()

		_, err = file.WriteString(string(jsonData))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			continue
		}

		puzzlesCreated += 1
	}
}
