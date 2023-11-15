package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	AuthUser(username string, password string) (role string)
	CreateEmployer(username string, password string, cash int) (err error)
	CreateWorker(username string, password string, weight int, wage int, drinks bool) (err error)
	GetWorker(username string) (worker Worker, err error)
	GetWorkerTasks(username string) (tasks []Completed)
	GetEmployer(username string) (employer Employer, workers []Worker, err error)
	GetEmployerTasks(username string) (tasks []Task)
}

type PostgreSQL struct {
	conn *pgxpool.Pool
}

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

func (p PostgreSQL) CreateEmployer(username string, password string, cash int) (err error) {

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

func (p PostgreSQL) CreateWorker(username string, password string, weight, wage int, drinks bool) (err error) {

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
	//worker, err = pgx.CollectOneRow[Worker](rows, pgx.RowToStructByNameLax[Worker])
	rows.Scan(&worker.userid, &worker.wage, &worker.fatigue, &worker.weight, &worker.drinks)
	rows.Close()
	return
}

func (p PostgreSQL) GetWorkerTasks(username string) (tasks []Completed) {
	return
}

func (p PostgreSQL) GetEmployer(username string) (employer Employer, workers []Worker, err error) {
	return
}

func (p PostgreSQL) GetEmployerTasks(username string) (tasks []Task) {
	return
}
