package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	AuthUser(username string, password string) (role string)
	CreateUser(username string, password string, role string) (err error)
	GetWorker(username string) (worker Worker, err error)
	GetWorkerTasks(username string) (tasks []Completed)
	GetEmployer(username string) (employer Employer, workers []Worker, err error)
	GetEmployerTasks(username string) (tasks []Task)
}

type PostgreSQL struct {
	conn *pgxpool.Pool
}

func (p PostgreSQL) AuthUser(username string, password string) (role string) {
	return ""
}

func (p PostgreSQL) CreateUser(username string, password string, role string) (err error) {
	return
}

func (p PostgreSQL) GetWorker(username string) (worker Worker, err error) {
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
