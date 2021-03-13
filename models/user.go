package models

import (
	"encoding/json"
)

// User models
type User struct {
	Name     string `firestore:"name" json:"name"`
	Photo    string `firestore:"photo" json:"photo"`
	Bio      string `firestore:"bio" json:"bio"`
	Phone    string `firestore:"phone" json:"phone"`
	Email    string `firestore:"email" json:"email"`
	Password string `firestore:"password" json:"password"`
}

// ToJSON Convert User to JSON
func (model User) ToJSON() ([]byte, error) {
	b, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	return b, nil
}
