package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/umenorin/rinha-backend-2023/pkg/domain"
	"github.com/umenorin/rinha-backend-2023/pkg/util"
)

type RouterManager struct {
	ConnectionDatabase *pgx.Conn
}

func (r *RouterManager) ExecRouter(port string) error {
	http.HandleFunc("GET /contagem-pessoas", r.CountPersons)
	http.HandleFunc("GET /pessoas/{id}", r.GetPersonById)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r RouterManager) CountPersons(w http.ResponseWriter, req *http.Request) {
	rows, err := r.ConnectionDatabase.Query(context.Background(), "SELECT COUNT(id)FROM person")
	util.CheckErrorQuery(err)
	defer rows.Close()
	for rows.Next() {
		var totalPeople int

		err = rows.Scan(&totalPeople)
		if err != nil {
			fmt.Printf("Scan error: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%d", totalPeople)
	}

}

func (r RouterManager) GetPersonById(w http.ResponseWriter, req *http.Request) {
	id, _ := strings.CutPrefix(req.URL.String(), "/pessoas/")
	rows, _ := r.ConnectionDatabase.Query(context.Background(), `
	SELECT p.id,p.name,p.nickname,p.birthdate, l.name FROM person as p
	LEFT JOIN stack as s ON p.id = s.person_id
  LEFT JOIN language as l ON s.language_id = l.id
  WHERE p.id =$1`, id)

	defer rows.Close()

	var person domain.Person
	for rows.Next() {
		var id, name, nickname, language string
		var birthdate time.Time
		rows.Scan(&id, &name, &nickname, &birthdate, &language)

		if language != "" {
			person.Language = append(person.Language, language)
		}

		person = domain.Person{Id: id, Name: name, Nickname: nickname, BirthDate: birthdate.Format("02-01-2006"), Language: person.Language}

		fmt.Printf("birthDate: %v\n", birthdate)
	}

	if !rows.CommandTag().Select() {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "ERROR 404 - USER NOTE FOUND")
		return
	}

	data, err := json.Marshal(person)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", data)
}

func (r RouterManager) GetPersonBySearch(w http.ResponseWriter, req *http.Request) {

}
