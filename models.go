package main

import (
	"math/rand"
	"time"
)

var Vowels = []string{
	"A",
	"E",
	"I",
	"O",
	"U",
}

var Consonants = []string{
	"B", "C", "D", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "V", "W", "X", "Y", "Z",
}

type Puzzle struct {
	Letters   []string
	Questions []Question
}

type Question struct {
	Hint   string
	Answer string
	Author string
}

func TakeRandomSubset[T any](slice []T, length int) []T {
	if length > len(slice) {
		length = len(slice)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	return slice[:length]
}

func (p Puzzle) WithRandomSubsetOfQuestions(length int) Puzzle {
	return Puzzle{
		Letters:   p.Letters,
		Questions: TakeRandomSubset(p.Questions, length),
	}
}

func (p Puzzle) FilterQuestions(f func(q Question) bool) Puzzle {
	newPuzzle := Puzzle{
		Letters:   p.Letters,
		Questions: make([]Question, 0),
	}

	for _, q := range p.Questions {
		if f(q) {
			newPuzzle.Questions = append(newPuzzle.Questions, q)
		}
	}

	return newPuzzle
}

func (p Puzzle) FilterDuplicateQuestionAnswers() Puzzle {
	answerSet := make(map[string]interface{}, 0)

	return p.FilterQuestions(func(q Question) bool {
		_, ok := answerSet[q.Answer]
		answerSet[q.Answer] = struct{}{}
		return !ok
	})
}
