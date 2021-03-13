package api

import (
	"errors"
	"net/http"
	"time"

	models "github.com/ngosangns/devchallenges-my-unsplash-api/models"
)

// Signup handler
func Signup(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	setHeader(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var rec models.User
		rec.Email = r.FormValue("email")
		rec.Password = r.FormValue("password")

		// Validate
		emailPattern := "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])"
		if !regEx(rec.Email, emailPattern) {
			printErr(w, errors.New("Email doesn't match pattern"), "")
			return
		}

		// Connect DB
		client, ctx, err := connectDb()
		defer client.Close()
		if err != nil {
			printErr(w, err, "Error while connecting to database")
			return
		}
		// Check user exist
		q := client.Collection("users").Where("email", "==", rec.Email)
		iter1 := q.Documents(ctx)
		defer iter1.Stop() // add this line to ensure resources cleaned up
		arr, _ := iter1.GetAll()
		if len(arr) > 0 {
			printErr(w, errors.New("Account already exists"), "")
			return
		}
		// Add new user
		_, _, err = client.Collection("users").Add(ctx, rec)
		if err != nil {
			printErr(w, err, "Error")
			return
		}
		// Create jwt token
		jwt, key := createToken(rec)
		now := time.Now().Add(time.Hour * 24 * 7)
		token := models.Token{
			JWT:     jwt,
			Key:     key,
			Expired: now,
			Email:   rec.Email,
		}
		// Add new token
		_, _, err = client.Collection("token").Add(ctx, token)
		if err != nil {
			printErr(w, err, "Error")
			return
		}

		// Print response
		printRes(w, map[string]interface{}{
			"token": token.JWT,
		})
	}
}
