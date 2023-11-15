CREATE TABLE users (
	userID SERIAL PRIMARY KEY,
  	username VARCHAR(200) NOT NULL,
  	password VARCHAR(200) NOT NULL,
    role VARCHAR(200) NOT NULL,
	UNIQUE(username)
);

CREATE TABLE employers (
	userID SERIAL REFERENCES users (userID),
	cash int NOT NULL,
	UNIQUE(userID)
);

CREATE TABLE workers (
	userID SERIAL REFERENCES users (userID),
	wage int NOT NULL,
    fatigue int NOT NULL,
	weight int NOT NULL,
    drinks boolean NOT NULL,
	UNIQUE(userID)
);

CREATE TABLE tasks (
	taskID SERIAL PRIMARY KEY,
    weight int,
	completed boolean
);

CREATE TABLE completed (
	taskID int NOT NULL,
	weight int NOT NULL,
	employer int NOT NULL,
    worker int NOT NULL
);