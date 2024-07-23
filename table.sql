CREATE TABLE clue (
  id SERIAL PRIMARY KEY,
  clueText VARCHAR (50) NOT NULL,
  answer VARCHAR (50) NOT NULL,
  author VARCHAR (50) NOT NULL,
  UNIQUE (clueText, answer)
);

CREATE TABLE letter (
  character CHAR(1) PRIMARY KEY 
);

CREATE TABLE clueLetterMapping (
    id SERIAL PRIMARY KEY,
    clue_id INTEGER NOT NULL,
    letter CHAR(1) NOT NULL,
    FOREIGN KEY (clue_id) REFERENCES clue(id),
    FOREIGN KEY (letter) REFERENCES letter(character),
    UNIQUE (clue_id, letter)
);

CREATE TABLE word (
  word VARCHAR(50) PRIMARY KEY
);

INSERT INTO letter (character) VALUES 
('A'), ('B'), ('C'), ('D'), ('E'), ('F'), ('G'), ('H'), ('I'), ('J'), 
('K'), ('L'), ('M'), ('N'), ('O'), ('P'), ('Q'), ('R'), ('S'), ('T'), 
('U'), ('V'), ('W'), ('X'), ('Y'), ('Z');
