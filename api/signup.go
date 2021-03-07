package api

import (
	"encoding/json"
	"net/http"
	"time"

	models "github.com/ngosangns/devchallenges-my-unsplash-api/models"
)

// Signup handler
func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var rec models.User

		rec.Email = r.FormValue("email")
		rec.Password = r.FormValue("password")

		// Connect DB
		client, ctx, err := connectDb()
		defer client.Close()
		if err != nil {
			printErr(w, err)
			return
		}
		// Add new user
		_, _, err = client.Collection("users").Add(ctx, rec)
		if err != nil {
			printErr(w, err)
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
			printErr(w, err)
			return
		}

		// Print response
		var res models.Res = models.Res{
			Status: true,
			Message: map[string]interface{}{
				"token": token.JWT,
			},
		}
		b, err := json.Marshal(res)
		if err != nil {
			printErr(w, err)
			return
		}
		w.Write(b)
	}
}
