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
	http.HandleFunc("GET /pessoas", r.GetPersonBySearch)
	http.HandleFunc("POST /pessoas", r.PostPeople)
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
	w.Header().Set("Content-Type", "application/json")
	id, _ := strings.CutPrefix(req.URL.String(), "/pessoas/")

	var person domain.Person

	person = FindById(r.ConnectionDatabase, id)
	if person.Nickname == "" {
		w.WriteHeader(http.StatusNotFound)
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
	w.Header().Set("Content-Type", "application/json")
	term := req.URL.Query().Get("t")
	if len(term) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "BAD REQUEST")
		return

	}
	term = "%" + term + "%"
	fmt.Printf("term: %v\n", term)
	rows, err := r.ConnectionDatabase.Query(context.Background(), `
	select DISTINCT p.id
	from person as p 
	left join stack as s on p.id = s.person_id
	left join language as l on l.id = s.language_id
	where LOWER( p.name) like LOWER($1) or LOWER(p.nickname) like LOWER($1) or LOWER(l.name) like LOWER($1)
	limit 50`, term)
	util.CheckErrorQuery(err)

	var person []domain.Person
	myIds := []string{}
	for rows.Next() {
		var id string
		rows.Scan(&id)

		fmt.Printf("id: %v\n", id)
		myIds = append(myIds, id)
	}
	rows.Close()
	for _, id := range myIds {

		people := FindById(r.ConnectionDatabase, id)
		person = append(person, people)
	}
	data, _ := json.Marshal(person)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", data)

}

func (r RouterManager) PostPeople(w http.ResponseWriter, req *http.Request) {

	var p domain.Person

	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(strings.ReplaceAll(p.Name, " ", "")) == 0 ||
		len(strings.ReplaceAll(p.Nickname, " ", "")) == 0 ||
		CheckNicknameAlreadyExist(r.ConnectionDatabase, p.Nickname) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	AddNewStackLanguageInDb(r.ConnectionDatabase, p)
}

func AddNewStackLanguageInDb(conn *pgx.Conn, people domain.Person) {
	conn.Exec(context.Background(), `
			INSERT INTO person (name,birthdate,nickname) VALUES ($1,$2,$3)
		`, people.Name, people.BirthDate, people.Nickname)

	if len(people.Language) > 0 {
		for _, lang := range people.Language {
			conn.Exec(context.Background(), `
			INSERT INTO language (name) VALUES ($1)
		`, lang)

		}
	}

	userId := getIdByNickName(conn, people.Nickname)
	fmt.Printf("userId: %v\n", userId)
	stacksId := getLanguagesIdByName(conn, people.Language)
	fmt.Printf("stacksId: %v\n", stacksId)
	for _, langId := range stacksId {
		conn.Exec(context.Background(), `
			INSERT INTO stack (person_id,language_id) VALUES ($1,$2)
		`, userId, langId)
	}
}

func CheckNicknameAlreadyExist(conn *pgx.Conn, nickname string) bool {
	rows, err := conn.Query(context.Background(), `
	SELECT p.id FROM person as p
  WHERE p.nickname =$1`, nickname)
	util.CheckErrorQuery(err)
	defer rows.Close()
	for rows.Next() {
		var id string
		rows.Scan(&id)
		if len(id) != 0 {
			return true
		}
	}

	return false
}

func getIdByNickName(conn *pgx.Conn, nickname string) string {
	rows, err := conn.Query(context.Background(), `
	SELECT p.id FROM person as p
  WHERE p.nickname =$1`, nickname)
	util.CheckErrorQuery(err)
	defer rows.Close()
	fmt.Printf("rows: %v\n", rows)
	for rows.Next() {
		var id string
		fmt.Printf("id: %v\n", id)
		rows.Scan(&id)
		if len(id) > 0 {
			return id
		}
	}

	return ""
}

func getLanguagesIdByName(conn *pgx.Conn, languages []string) []string {
	languagesId := []string{}
	for _, lang := range languages {
		rows, err := conn.Query(context.Background(), `
	SELECT l.id FROM language as l
  WHERE l.name =$1
	LIMIT 1
	`, lang)
		util.CheckErrorQuery(err)
		defer rows.Close()

		for rows.Next() {
			var id string
			rows.Scan(&id)
			if len(id) > 0 {
				languagesId = append(languagesId, id)
			}
		}
	}

	return languagesId
}
func FindById(conn *pgx.Conn, id string) domain.Person {
	rows, err := conn.Query(context.Background(), `
	SELECT p.id,p.name,p.nickname,p.birthdate, l.name FROM person as p
	LEFT JOIN stack as s ON p.id = s.person_id
  LEFT JOIN language as l ON s.language_id = l.id
  WHERE p.id =$1`, id)
	util.CheckErrorQuery(err)
	defer rows.Close()

	var person domain.Person
	for rows.Next() {
		var id, name, nickname, language string
		var birthdate time.Time
		rows.Scan(&id, &name, &nickname, &birthdate, &language)
		fmt.Printf("name: %v\n", name)
		if language != "" {
			person.Language = append(person.Language, language)
		}

		person = domain.Person{Id: id, Name: name, Nickname: nickname, BirthDate: birthdate.Format("2006-01-02"), Language: person.Language}

		fmt.Printf("birthDate: %v\n", birthdate)
	}
	fmt.Println(person)
	return person
}
