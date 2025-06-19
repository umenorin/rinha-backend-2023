package route

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type MyHandler struct {
	ConnectionDatabase *pgx.Conn
}

func (m MyHandler) GetPersons(w http.ResponseWriter, req *http.Request) {
	rows, err := m.ConnectionDatabase.Query(context.Background(), "select * from person")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		var birthDate time.Time
		var nickName string

		err = rows.Scan(&id, &name, &birthDate, &nickName)
		if err != nil {
			fmt.Printf("Scan error: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w,"%q\t | %q\t |%v\t |%q\t\n", id, name, birthDate, nickName)
	}

}
