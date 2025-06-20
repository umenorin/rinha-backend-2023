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

func (m *RouterManager) ExecRouter(port string) error {
	http.HandleFunc("GET /contagem-pessoas", m.CountPersons)
	http.HandleFunc("GET /pessoas/{id}", m.GetPersonById)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		return err
	}
	return nil
}

func (m RouterManager) CountPersons(w http.ResponseWriter, req *http.Request) {
	rows, err := m.ConnectionDatabase.Query(context.Background(), "SELECT COUNT(id)FROM person")
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
func (m RouterManager) GetPersonById(w http.ResponseWriter, req *http.Request) {
	path := req.URL
	id, _ := strings.CutPrefix(path.String(), "/pessoas/")
	rows, err := m.ConnectionDatabase.Query(context.Background(), "SELECT id,name,nickname,birthdate FROM person where person.id=$1", id)
	util.CheckErrorQuery(err)
	defer rows.Close()
	var person domain.Person
	for rows.Next() {
		var id, name, nickname string
		var birthdate time.Time
		err = rows.Scan(&id, &name, &nickname, &birthdate)
		if err != nil {
			fmt.Printf("Scan error: %v", err)
			return
		}
		person = domain.Person{Id: id, Name: name, Nickname: nickname, BirthDate: birthdate}
	}

	rows, err = m.ConnectionDatabase.Query(context.Background(), "SELECT l.name FROM language as l INNER JOIN stack as s on l.id = s.language_id INNER JOIN person as p on s.person_id=p.id where p.id=$1;", id)
	util.CheckErrorQuery(err)
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			fmt.Printf("Scan error: %v", err)
			return
		}
		person.Language = append(person.Language, name)
	}

	data, err := json.MarshalIndent(person, "", " ")
	if err != nil {

		log.Fatalf("JSON marshaling failed: %s", err)

	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", data)

}
