CREATE TABLE users (
	userID SERIAL PRIMARY KEY,
  	username VARCHAR(200) NOT NULL,
  	password VARCHAR(200) NOT NULL,
    role VARCHAR(200) NOT NULL
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
    drinks boolean,
	UNIQUE(userID)
);

CREATE TABLE tasks (
	taskID SERIAL PRIMARY KEY,
    employer SERIAL REFERENCES users (userID),
    weight int,
	completed boolean,
    UNIQUE(taskID)
);

CREATE TABLE taskTeams (
	taskID SERIAL REFERENCES tasks (taskID),
    worker SERIAL REFERENCES users (userID)
);