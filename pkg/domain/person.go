package domain

import "time"

type Person struct {
	Id        string    `json: "id"`
	Name      string    `json:"nome"`
	BirthDate time.Time `json:"nascimento"`
	Nickname  string    `json:"apelido"`
	Language  []string  `json:"stack"`
}
