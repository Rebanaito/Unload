package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	// Connecting to the PostgreSQL database. TODO: move credentials to environment file
	conn, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5432/unload")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating a pgx pool: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Starting the server on port 8080. TODO: make port part of env?
	server := NewAPIServer(":8080", conn)
	server.Run()
}
