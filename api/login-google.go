package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ngosangns/devchallenges-my-unsplash-api/models"
)

const googleClientID = "617831923199-ha6054jhlqqkrioohv5fioo5m5f10iki.apps.googleusercontent.com"
const googleClientSecret = "wtgMIEiAt5UGKjg3BiBNCIf5"
const googleRedirectURL = "https://ngosangns-authapp.web.app/login-google"
const googleLinkGetToken = "https://accounts.google.com/o/oauth2/token"
const googleLinkGetUserInfo = "https://www.googleapis.com/oauth2/v1/userinfo?access_token="
const googleGrantType = "authorization_code"

// LoginGoogle handler
func LoginGoogle(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	setHeader(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	// Get URL param "code"
	keys, ok := r.URL.Query()["code"]
	if !ok || len(keys[0]) < 1 {
		printErr(w, errors.New("URL param 'code' is missing"), "")
		return
	}
	key := keys[0]

	// Get Google API access token
	accessToken := getToken(key)
	// Get user info in Google Account
	var userInfo map[string]interface{}
	userInfo = getUserInfo(accessToken)

	// Connect DB
	client, ctx, err := connectDb()
	defer client.Close()
	if err != nil {
		printErr(w, err, "Error while connecting to database")
		return
	}
	// Get a record
	q := client.Collection("users").Where("email", "==", userInfo["email"])
	iter1 := q.Documents(ctx)
	defer iter1.Stop() // add this line to ensure resources cleaned up
	arr, _ := iter1.GetAll()
	// If account exist
	if len(arr) > 0 {
		rec := models.User{
			Email: fmt.Sprintf("%v", userInfo["email"]),
		}
		// Get token
		q = client.Collection("token").Where("email", "==", userInfo["email"])
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
			return
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
			return
		}
	} else { // If account doesn't exists
		rec := models.User{
			Email: fmt.Sprintf("%v", userInfo["email"]),
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
		return
	}
}

func getToken(code string) string {
	data := url.Values{
		"client_id":     {googleClientID},
		"client_secret": {googleClientSecret},
		"redirect_uri":  {googleRedirectURL},
		"code":          {code},
		"grant_type":    {googleGrantType},
	}
	resp, err := http.PostForm(googleLinkGetToken, data)

	if err != nil {
		return ""
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	accessToken := fmt.Sprintf("%v", res["access_token"])
	accessToken = strings.ReplaceAll(accessToken, "\"", "")
	return accessToken
}

func getUserInfo(accessToken string) map[string]interface{} {
	link := googleLinkGetUserInfo + accessToken
	resp, err := http.Get(link)

	if err != nil {
		log.Println("Error while requesting")
		return map[string]interface{}{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while requesting")
		return map[string]interface{}{}
	}

	r := bytes.NewReader(body)
	var res map[string]interface{}
	json.NewDecoder(r).Decode(&res)
	return res
}
