package model

import "github.com/google/uuid"

type Users struct {
	Data []User `json:"data"`
}

type User struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
}

type Envs struct {
	Data []Environment `json:"data"`
}

type Environment struct {
	ID 		uuid.UUID `json:"id"`
	Title 	string 	  `json:"title"`
	Text 	string 	  `json:"text"`
}

type Good struct {
	ID 		uuid.UUID `json:"id"`
	Name 	string 	  `json:"name"`
	Price 	int 	  `json:"price"` // for now int
	Qty 	int 	  `json:"qty"`
}


