package models

import time "time"

// Token models
type Token struct {
	Email   string    `firestore:"email"`
	JWT     string    `firestore:"jwt"`
	Key     string    `firestore:"key"`
	Expired time.Time `firestore:"expired"`
}
