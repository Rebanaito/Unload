package main

type Employer struct {
	userid int
	cash   int
}

type Worker struct {
	userid  int `db:"userid"`
	wage    int
	fatigue int
	drinks  bool
	weight  int
}

type Task struct {
	taskID int
	weight int
}

type Completed struct {
	taskID   int
	weight   int
	employer int
	userID   int
}
