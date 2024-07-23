package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	connection *sql.DB
}

func GetDatabase(host, user, name, password string, port int) Database {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, name, password)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return Database{
		connection: db,
	}
}

func (db *Database) addClue(clue, answer, author string) (int, error) {
	insertStatement := "INSERT INTO clue (clueText, answer, author) VALUES ($1, $2, $3)"

	_, err := db.connection.Exec(insertStatement, clue, answer, author)

	if err != nil {
		return -1, err
	}

	var id int

	if err := db.connection.QueryRow("SELECT id FROM clue WHERE clueText = $1 AND answer = $2",
		clue, answer).Scan(&id); err != nil {
		return -1, fmt.Errorf("%s", err.Error())
	}

	return id, nil
}

func (db *Database) addLetter(id int, letter rune) error {
	insertStatement := "INSERT INTO clueLetterMapping (clue_id, letter) VALUES ($1, $2)"

	_, err := db.connection.Exec(insertStatement, id, string(letter))
	return err
}

func (db *Database) addWord(word string) error {
	insertStatement := "INSERT INTO word(word) VALUES ($1)"

	_, err := db.connection.Exec(insertStatement, word)
	return err
}

// todo why can't I use a []string???
func (db *Database) queryWordsWithLetters(a, b, c, d, e, f, g string) ([]Question, error) {
	query := `
		-- Create a temporary table with the allowed letters
		WITH allowed_letters AS (
			SELECT $1 AS character
			UNION ALL SELECT $2
			UNION ALL SELECT $3
			UNION ALL SELECT $4
			UNION ALL SELECT $5
			UNION ALL SELECT $6
			UNION ALL SELECT $7
		),
		-- Find all clues where the letters are within the allowed subset
		clues_with_allowed_letters AS (
			SELECT clm.clue_id
			FROM clueLetterMapping clm
			JOIN allowed_letters al ON clm.letter = al.character
			GROUP BY clm.clue_id
			HAVING COUNT(clm.letter) = (
				SELECT COUNT(*)
				FROM clueLetterMapping clm2
				WHERE clm2.clue_id = clm.clue_id
			)
		)
		-- Select clues that match the above criteria
		SELECT c.clueText, c.answer, c.author
		FROM clue c
		JOIN clues_with_allowed_letters cal ON c.id = cal.clue_id
		-- Filter to only include answers in the english dictionary
		JOIN word ON c.answer = word.word
		-- Only allow words of a certain length
		WHERE LENGTH(c.answer) >= 4;
	`

	rows, err := db.connection.Query(query, a, b, c, d, e, f, g)
	if err != nil {
		return make([]Question, 0), err
	}
	defer rows.Close()

	var results []Question

	for rows.Next() {
		var value Question
		if err := rows.Scan(&value.Hint, &value.Answer, &value.Author); err != nil {
			return make([]Question, 0), err
		}
		results = append(results, value)
	}

	if err := rows.Err(); err != nil {
		return make([]Question, 0), err
	}

	return results, nil
}
