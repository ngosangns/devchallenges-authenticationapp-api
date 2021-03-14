package api

import (
	"errors"
	"fmt"
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
		// Get token from Authorization header
		token := r.Header.Get("Authorization")
		// Check token pattern
		authPattern := `Bearer .+`
		if !regEx(token, authPattern) {
			printErr(w, errors.New("Authorization is required"), "")
			return
		}
		token = token[7:]

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
					printErr(w, err, "Error")
					return
				}
				rec.Password = "" // Hide password
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

	if r.Method == "POST" {
		// Get token from Authorization header
		token := r.Header.Get("Authorization")
		// Check token pattern
		authPattern := `Bearer .+`
		if !regEx(token, authPattern) {
			printErr(w, errors.New("Authorization is required"), "")
			return
		}
		token = token[7:]

		// Get upload info
		var rec models.User
		rec.Email = r.FormValue("email")
		rec.Password = r.FormValue("password")
		rec.Photo = r.FormValue("photo")
		rec.Bio = r.FormValue("bio")
		rec.Phone = r.FormValue("phone")
		rec.Name = r.FormValue("name")

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
			// Get user record
			q := client.Collection("users").Where("email", "==", email)
			iter1 := q.Documents(ctx)
			defer iter1.Stop() // add this line to ensure resources cleaned up
			arr, _ := iter1.GetAll()
			// If account exist
			if len(arr) > 0 {
				// Check upload password
				if rec.Password == "" {
					recordPassword, _ := arr[0].DataAt("password")
					rec.Password = fmt.Sprintf("%v", recordPassword)
				}
				// Update account info
				_, err := client.Collection("users").Doc(arr[0].Ref.ID).Set(ctx, rec)
				if err != nil {
					printErr(w, err, "Error")
					return
				}
				printRes(w, "Update successful")
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
