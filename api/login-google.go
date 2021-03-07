package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const googleClientID = "617831923199-ha6054jhlqqkrioohv5fioo5m5f10iki.apps.googleusercontent.com"
const googleClientSecret = "wtgMIEiAt5UGKjg3BiBNCIf5"
const googleRedirectURL = "http://localhost:8080/api/login-google"
const googleLinkGetToken = "https://accounts.google.com/o/oauth2/token"
const googleLinkGetUserInfo = "https://www.googleapis.com/oauth2/v1/userinfo?access_token="
const googleGrantType = "authorization_code"

// LoginGoogle handler
func LoginGoogle(w http.ResponseWriter, r *http.Request) {
	// Get URL param "code"
	keys, ok := r.URL.Query()["code"]
	if !ok || len(keys[0]) < 1 {
		fmt.Fprintf(w, "URL param 'code' is missing")
	}
	key := keys[0]
	accessToken := getToken(key)

	var userInfo map[string]interface{}
	userInfo = getUserInfo(accessToken)

	fmt.Fprintf(w, fmt.Sprintf("%v", map[string]interface{}{
		"id":    userInfo["id"],
		"name":  userInfo["name"],
		"email": userInfo["email"],
	}))
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
