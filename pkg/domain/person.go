package domain

type Person struct {
	Id        string   `json: "id"`
	Name      string   `json:"nome"`
	BirthDate string   `json:"nascimento"`
	Nickname  string   `json:"apelido"`
	Language  []string `json:"stack"`
}
