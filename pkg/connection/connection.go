package connection

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)
// TODO 
//fazer essa parte ser chamada pelo main.go sem ter o problema de espera
//possivel soluçao coccurency e chanels
func StartConnection() (context.Context, *pgx.Conn) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	fmt.Println("Starting connection with Postgres Db")
	conn, err := pgx.Connect(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	err = http.ListenAndServe(":8000", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")

	} else if err != nil {
		fmt.Printf("error startin server: %s\n", err)
		os.Exit(1)
	}
	return ctx, conn
}
