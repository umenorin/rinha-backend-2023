package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/umenorin/rinha-backend-2023/pkg/router"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASS")
	HOST := os.Getenv("DB_HOST")
	PORT := os.Getenv("DB_PORT")
	DBNAME := os.Getenv("DB_DBNAME")
	connection := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", USER, PASS, HOST, PORT, DBNAME)

	ctx := context.Background()
	fmt.Printf("Starting connection with Postgres Db in connection %v\n",connection)
	conn, err := pgx.Connect(ctx, connection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	PORTAPI := os.Getenv("API_PORT")
	route := router.RouterManager{ConnectionDatabase: conn}
	route.ExecRouter(fmt.Sprintf(":%v",PORTAPI))

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")

	} else if err != nil {
		fmt.Printf("error startin server: %s\n", err)
		os.Exit(1)
	}

}
