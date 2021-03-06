package main

import (
	"net/http"

	api "github.com/ngosangns/devchallenges-my-unsplash-api/api"
)

func main() {
	http.HandleFunc("/api/date", api.Date)
	http.HandleFunc("/api/login-google", api.LoginGoogle)
	http.HandleFunc("/api/db", api.Db)
	http.ListenAndServe(":8080", nil)
}
