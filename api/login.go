package api

import (
	"errors"
	"net/http"
	"time"

	models "github.com/ngosangns/devchallenges-my-unsplash-api/models"
)

// Login handler
func Login(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	setHeader(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var rec models.User
		rec.Email = r.FormValue("email")
		rec.Password = r.Form.Get("password")

		// Validate
		emailPattern := `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
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
		// Get a record
		q := client.Collection("users").Where("email", "==", rec.Email)
		iter1 := q.Documents(ctx)
		defer iter1.Stop() // add this line to ensure resources cleaned up
		arr, _ := iter1.GetAll()
		// If account exist
		if len(arr) > 0 {
			password, _ := arr[0].DataAt("password")
			if password == rec.Password {
				// Get token
				q = client.Collection("token").Where("email", "==", rec.Email)
				iter2 := q.Documents(ctx)
				defer iter2.Stop()
				arr, _ = iter2.GetAll()
				// If account has token then write the token to response
				if len(arr) > 0 {
					jwt, err := arr[0].DataAt("jwt")
					if err != nil {
						printErr(w, err, "Error")
						return
					}
					// Write token to response
					printRes(w, map[string]interface{}{
						"token": jwt,
					})
				} else { // If token doesn't exist then create a new one
					// Create jwt token
					jwt, key := createToken(rec)
					now := time.Now().Add(time.Hour * 24 * 7)
					token := models.Token{
						JWT:     jwt,
						Key:     key,
						Expired: now,
						Email:   rec.Email,
					}
					// Add new token to database
					_, _, err = client.Collection("token").Add(ctx, token)
					if err != nil {
						printErr(w, err, "Error")
						return
					}
					// Write token to response
					printRes(w, map[string]interface{}{
						"token": token.JWT,
					})
				}
			} else {
				printErr(w, errors.New("Wrong password"), "")
			}
			return
		}
		printErr(w, errors.New("Account doesn't exist"), "")
	}
}
