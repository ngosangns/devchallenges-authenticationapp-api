package api

import (
	"errors"
	"net/http"

	"github.com/ngosangns/devchallenges-my-unsplash-api/models"
)

// User handler
func User(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	setHeader(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "GET" {
		// Get URL param "token"
		keys, ok := r.URL.Query()["token"]
		if !ok || len(keys[0]) < 1 {
			printErr(w, errors.New("URL param 'code' is missing"), "")
			return
		}
		token := keys[0]

		// Connect DB
		client, ctx, err := connectDb()
		defer client.Close()
		if err != nil {
			printErr(w, err, "Error while connecting to database")
			return
		}
		// Check token
		q := client.Collection("token").Where("jwt", "==", token)
		iter1 := q.Documents(ctx)
		defer iter1.Stop() // add this line to ensure resources cleaned up
		arr, _ := iter1.GetAll()
		// If token exist
		if len(arr) > 0 {
			email, err := arr[0].DataAt("email")
			if err != nil {
				printErr(w, err, "Error")
				return
			}
			// Get user info
			q := client.Collection("users").Where("email", "==", email)
			iter1 := q.Documents(ctx)
			defer iter1.Stop() // add this line to ensure resources cleaned up
			arr, _ := iter1.GetAll()
			// If account exist
			if len(arr) > 0 {
				var rec models.User
				err = arr[0].DataTo(&rec)
				if err != nil {
					printErr(w, err, "")
					return
				}
				printRes(w, rec)
				return
			} else { // If account doesn't exists
				printErr(w, errors.New("Wrong info"), "")
				return
			}
		} else { // If token doesn't exists
			printErr(w, errors.New("Wrong info"), "")
			return
		}
	}

}
