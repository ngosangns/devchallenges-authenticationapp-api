package main

import (
	"net/http"

	api "github.com/ngosangns/devchallenges-my-unsplash-api/api"
)

func main() {
	http.HandleFunc("/api/date", api.Date)
	http.HandleFunc("/api/db", api.Db)
	http.HandleFunc("/api/signup", api.Signup)
	http.HandleFunc("/api/login", api.Login)
	http.HandleFunc("/api/login-google", api.LoginGoogle)
	http.HandleFunc("/api/user", api.User)
	http.ListenAndServe(":8080", nil)
}
