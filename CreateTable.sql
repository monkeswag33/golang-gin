CREATE TABLE mytable (
	id SERIAL PRIMARY KEY,
	firstname TEXT NOT NULL,
	lastname TEXT NOT NULL
);

INSERT INTO mytable (firstname, lastname) VALUES ('User1First', 'User1Last'), ('User2First', 'User2Last');