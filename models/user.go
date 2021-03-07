package models

import (
	"encoding/json"
)

// User models
type User struct {
	Name     string `firestore:"name"`
	Photo    string `firestore:"photo"`
	Bio      string `firestore:"bio"`
	Phone    string `firestore:"phone"`
	Email    string `firestore:"email"`
	Password string `firestore:"password"`
}

// ToJSON Convert User to JSON
func (model User) ToJSON() ([]byte, error) {
	b, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	return b, nil
}
