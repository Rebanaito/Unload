package main

type Employer struct {
	userid int
	cash   int
}

type Worker struct {
	userid  int
	wage    int
	fatigue int
	drinks  bool
	weight  int
}

type Task struct {
	taskID    int
	weight    int
	completed bool
}

type Completed struct {
	taskID   int
	weight   int
	employer int
	worker   int
}
