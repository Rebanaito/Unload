package main

type Cache struct {
	users          map[Credentials]User
	employers      map[int]Employer
	workers        map[int]Worker
	completedTasks []Task
	availableTasks []Task
}

type User struct {
	userID   int
	username string
	password string
	role     string
}

type Employer struct {
	userID int
	cash   int
	tasks  []*Task
}

type Worker struct {
	userID  int
	wage    int
	fatigue int
	drinks  bool
	tasks   []*Task
}

type Task struct {
	taskID    int
	workers   []*Worker
	weight    int
	completed bool
}

type Credentials struct {
	username string
	password string
}
