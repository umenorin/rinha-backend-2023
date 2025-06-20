package domain

import "time"

type Person struct {
	Id        string    `json: "id"`
	Name      string    `json:"name"`
	BirthDate time.Time `json:"birthdate"`
	Nickname  string    `json:"nickname"`
}
