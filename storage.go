package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Interface that allows us to more easily implement the game with different storage
type Storage interface {
	AuthUser(username string, password string) (role string)
	CreateEmployer(username string, password string, cash int) (err error)
	CreateWorker(username string, password string, weight int, wage int, drinks bool) (err error)
	GetWorker(username string) (worker Worker, err error)
	GetWorkerByID(id string) (worker Worker)
	GetWorkerTasks(username string) (tasks []Completed)
	GetEmployer(username string) (employer Employer, workers []Worker, err error)
	GetAllWorkers() (workers []Worker)
	GetEmployerTasks() (tasks []Task)
	GetActiveTaskCount() (active int)
	AddTask(weight int)
	GetTask(taskID string) (task Task)
	AddMoney(username string, profit int)
	MarkComplete(task Task, username string, workers []Worker)
	UpdateWorkers(workers []Worker)
	RemoveMoney(username string, wageTotal int)
	UpdateTask(task Task, weight int)
}

// Current implementation of Storage
type PostgreSQL struct {
	conn *pgxpool.Pool
}

// Finds user with provided credentials, if found returns their role
func (p PostgreSQL) AuthUser(username string, password string) (role string) {
	query := fmt.Sprintf("SELECT role FROM users WHERE username='%s' AND password='%s'", username, password)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database: %v\n", err)
		os.Exit(1)
	}
	r, err := pgx.CollectOneRow[string](rows, pgx.RowTo)
	if err == nil && (r == "employer" || r == "worker") {
		role = r
	}
	return
}

// Create a new Employer user with provided credentials and money. Returns error if user already exists
// or if the funds are insufficient
func (p PostgreSQL) CreateEmployer(username string, password string, cash int) (err error) {

	// Check that the provided amount is valid
	if cash < 0 || cash > 100000 {
		return errors.New("invalid amount")
	}

	// Check for existing non-bankrupt employers
	rows, err := p.conn.Query(context.Background(), "SELECT COUNT(*) FROM employers WHERE cash > 0")
	if err != nil {
		return err
	}
	count, _ := pgx.CollectOneRow[int](rows, pgx.RowTo)
	if count != 0 {
		return errors.New("employer exists and is not bankrupt")
	}

	// Try adding a user, checking for duplicate username
	query := fmt.Sprintf("INSERT INTO users (username, password, role) VALUES ('%s', '%s', 'employer')", username, password)
	rows, err = p.conn.Query(context.Background(), query)
	if err != nil {
		return err
	}
	_, err = pgx.CollectOneRow[string](rows, pgx.RowTo)
	if err != nil && err.Error() != "no rows in result set" {
		return errors.New("username already exists")
	}

	// Fetching the ID of the newly created user
	query = fmt.Sprintf("SELECT userID FROM users WHERE username='%s'", username)
	rows, err = p.conn.Query(context.Background(), query)
	if err != nil {
		return err
	}
	userID, err := pgx.CollectOneRow[int](rows, pgx.RowTo)
	if err != nil {
		return errors.New("unexpected error, try again")
	}

	// Inserting into employers table
	query = fmt.Sprintf("INSERT INTO employers (userID, cash) VALUES ('%d', '%d')", userID, cash)
	rows, err = p.conn.Query(context.Background(), query)
	if err != nil {
		return err
	}
	_, err = pgx.CollectOneRow[string](rows, pgx.RowTo)
	if err != nil && err.Error() != "no rows in result set" {
		return errors.New("unexpected error")
	}
	return nil
}

// Create a new Worker user, with the provided credentials and parameters. Return error if values are invalid
func (p PostgreSQL) CreateWorker(username string, password string, weight, wage int, drinks bool) (err error) {

	// Checking for invalid parameters
	if weight < 5 || weight > 30 {
		return errors.New("invalid weight")
	}

	if wage < 10000 || weight > 30000 {
		return errors.New("invalid wage")
	}

	// Try adding a user, checking for duplicate username
	query := fmt.Sprintf("INSERT INTO users (username, password, role) VALUES ('%s', '%s', 'worker')", username, password)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return errors.New("unexpected error, try again")
	}
	_, err = pgx.CollectOneRow[string](rows, pgx.RowTo)
	if err != nil && err.Error() != "no rows in result set" {
		return errors.New("username already exists")
	}

	// Fetching the ID of the newly created user
	query = fmt.Sprintf("SELECT userID FROM users WHERE username='%s'", username)
	rows, err = p.conn.Query(context.Background(), query)
	if err != nil {
		return errors.New("unexpected error, try again")
	}
	userID, err := pgx.CollectOneRow[int](rows, pgx.RowTo)
	if err != nil {
		return errors.New("unexpected error, try again")
	}

	// Inserting into workers table
	query = fmt.Sprintf("INSERT INTO workers (userID, weight, wage, fatigue, drinks) VALUES ('%d', '%d', '%d', '0', '%v')", userID, weight, wage, drinks)
	rows, err = p.conn.Query(context.Background(), query)
	if err != nil {
		return errors.New("unexpected error, try again")
	}
	_, err = pgx.CollectOneRow[string](rows, pgx.RowTo)
	if err != nil && err.Error() != "no rows in result set" {
		return errors.New("unexpected error, try again")
	}
	return nil
}

// Get a worker object by username. TODO: make JWT carry id instead of username
func (p PostgreSQL) GetWorker(username string) (worker Worker, err error) {
	query := fmt.Sprintf("SELECT userID FROM users WHERE username='%s'", username)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	userID, err := pgx.CollectOneRow[int](rows, pgx.RowTo)
	if err != nil {
		return
	}
	query = fmt.Sprintf("SELECT * FROM workers WHERE userid='%d'", userID)
	rows, err = p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	rows.Next()
	rows.Scan(&worker.userid, &worker.wage, &worker.fatigue, &worker.weight, &worker.drinks)
	rows.Close()
	return
}

// Find the tasks completed by the worker
func (p PostgreSQL) GetWorkerTasks(username string) (tasks []Completed) {
	query := fmt.Sprintf("SELECT userID FROM users WHERE username='%s'", username)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	userID, err := pgx.CollectOneRow[int](rows, pgx.RowTo)
	if err != nil {
		return
	}
	query = fmt.Sprintf("SELECT * FROM completed WHERE worker='%d'", userID)
	rows, err = p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	rows.Next()
	for {
		var c Completed
		rows.Scan(&c.taskID, &c.weight, &c.employer, &c.worker)
		if err != nil {
			fmt.Println(err)
		}
		tasks = append(tasks, c)
		if !rows.Next() {
			break
		}
	}
	return
}

// Get the employer and all available workers
func (p PostgreSQL) GetEmployer(username string) (employer Employer, workers []Worker, err error) {
	query := fmt.Sprintf("SELECT userID FROM users WHERE username='%s'", username)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	employer.userid, err = pgx.CollectOneRow[int](rows, pgx.RowTo)
	if err != nil {
		return
	}
	query = fmt.Sprintf("SELECT cash FROM employers WHERE userid='%d'", employer.userid)
	rows, err = p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	employer.cash, err = pgx.CollectOneRow[int](rows, pgx.RowTo)
	if err != nil {
		return
	}
	workers = p.GetAllWorkers()
	return
}

func (p PostgreSQL) GetAllWorkers() (workers []Worker) {
	query := "SELECT * FROM workers"
	rows, _ := p.conn.Query(context.Background(), query)
	rows.Next()
	for {
		var w Worker
		rows.Scan(&w.userid, &w.wage, &w.fatigue, &w.weight, &w.drinks)
		workers = append(workers, w)
		if !rows.Next() {
			break
		}
	}
	return
}

// Get all available (not completed) tasks
func (p PostgreSQL) GetEmployerTasks() (tasks []Task) {
	query := "SELECT * FROM tasks WHERE completed='false'"
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	rows.Next()
	for {
		var t Task
		rows.Scan(&t.taskID, &t.weight, &t.completed)
		if err != nil {
			fmt.Println(err)
		}
		tasks = append(tasks, t)
		if !rows.Next() {
			break
		}
	}
	return
}

// Count the number of available tasks. Used by the task generating routine
func (p PostgreSQL) GetActiveTaskCount() (active int) {
	query := "SELECT * FROM tasks WHERE completed='false'"
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	for rows.Next() {
		active++
	}
	return
}

// Adds a new task. Used by the task generating routine
func (p PostgreSQL) AddTask(weight int) {
	query := fmt.Sprintf("INSERT INTO tasks (weight, completed) VALUES ('%d', 'false')", weight)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		fmt.Println(err)
	}
	rows.Close()
}

func (p PostgreSQL) GetTask(taskID string) (task Task) {
	query := fmt.Sprintf("SELECT * FROM tasks WHERE taskid='%s'", taskID)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	rows.Next()
	rows.Scan(&task.taskID, &task.weight, &task.completed)
	rows.Close()
	return
}

func (p PostgreSQL) GetWorkerByID(id string) (worker Worker) {
	query := fmt.Sprintf("SELECT * FROM workers WHERE userid='%s'", id)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	rows.Next()
	rows.Scan(&worker.userid, &worker.wage, &worker.fatigue, &worker.weight, &worker.drinks)
	rows.Close()
	return
}

// Adds money to the employer upon successful completion of the task
func (p PostgreSQL) AddMoney(username string, profit int) {
	query := fmt.Sprintf("SELECT userID FROM users WHERE username='%s'", username)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	userid, err := pgx.CollectOneRow[int](rows, pgx.RowTo)
	if err != nil {
		return
	}
	query = fmt.Sprintf("SELECT cash FROM employers WHERE userid='%d'", userid)
	rows, err = p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	cash, err := pgx.CollectOneRow[int](rows, pgx.RowTo)
	if err != nil {
		return
	}
	query = fmt.Sprintf("UPDATE employers SET cash='%d' WHERE userid='%d'", cash+profit, userid)
	p.conn.Query(context.Background(), query)
}

// Marks the task as complete and saves the data about the team composition (who worked on the task)
func (p PostgreSQL) MarkComplete(task Task, username string, workers []Worker) {
	query := fmt.Sprintf("SELECT userID FROM users WHERE username='%s'", username)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	employer, err := pgx.CollectOneRow[int](rows, pgx.RowTo)
	if err != nil {
		return
	}
	for _, worker := range workers {
		query := fmt.Sprintf("INSERT INTO completed (taskid, weight, employer, worker) VALUES ('%d', '%d', '%d', '%d')",
			task.taskID, task.weight, employer, worker.userid)
		p.conn.Query(context.Background(), query)
	}
	query = fmt.Sprintf("UPDATE tasks SET completed='true' WHERE taskid='%d'", task.taskID)
	p.conn.Query(context.Background(), query)
}

// Count the new fatigue values
func (p PostgreSQL) UpdateWorkers(workers []Worker) {
	for _, worker := range workers {
		if worker.fatigue >= 100 {
			continue
		}
		if worker.drinks {
			worker.fatigue += 30
		} else {
			worker.fatigue += 20
		}
		if worker.fatigue > 100 {
			worker.fatigue = 100
		}
		query := fmt.Sprintf("UPDATE workers SET fatigue='%d' WHERE userid='%d'", worker.fatigue, worker.userid)
		p.conn.Query(context.Background(), query)
	}
}

// Penalizes the employer after he fails the task
func (p PostgreSQL) RemoveMoney(username string, wageTotal int) {
	query := fmt.Sprintf("SELECT userID FROM users WHERE username='%s'", username)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	userid, err := pgx.CollectOneRow[int](rows, pgx.RowTo)
	if err != nil {
		return
	}
	query = fmt.Sprintf("SELECT cash FROM employers WHERE userid='%d'", userid)
	rows, err = p.conn.Query(context.Background(), query)
	if err != nil {
		return
	}
	cash, err := pgx.CollectOneRow[int](rows, pgx.RowTo)
	if err != nil {
		return
	}
	query = fmt.Sprintf("UPDATE employers SET cash='%d' WHERE userid='%d'", cash-(wageTotal+wageTotal/10), userid)
	p.conn.Query(context.Background(), query)
}

func (p PostgreSQL) UpdateTask(task Task, weight int) {
	query := fmt.Sprintf("UPDATE tasks SET weight='%d' WHERE taskid='%d'", task.weight-weight, task.taskID)
	p.conn.Query(context.Background(), query)
}
