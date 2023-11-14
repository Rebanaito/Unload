package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	AuthUser(username string, password string) (id int, role string)
	CreateUser(username string, password string, role string) (err error)
	GetWorker(id int) (worker Worker, err error)
	GetWorkerTasks(id int) (tasks []Task)
	GetEmployer(id int) (employer Employer, workers []Worker, err error)
	GetEmployerTasks(id int) (tasks []Task)
}

type PostgreSQL struct {
	conn *pgxpool.Pool
}

func (p PostgreSQL) AuthUser(username string, password string) (id int, role string) {
	return -1, ""
}

func (p PostgreSQL) CreateUser(username string, password string, role string) (err error) {
	return
}

func (p PostgreSQL) GetWorker(id int) (worker Worker, err error) {
	return
}

func (p PostgreSQL) GetWorkerTasks(id int) (tasks []Task) {
	return
}

func (p PostgreSQL) GetEmployer(id int) (employer Employer, workers []Worker, err error) {
	return
}

func (p PostgreSQL) GetEmployerTasks(id int) (tasks []Task) {
	return
}
