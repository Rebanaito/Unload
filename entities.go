package main

type Employer struct {
	cash int
}

type Worker struct {
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
	employer int
	weight   int
}
