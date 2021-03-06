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

const GOOGLE_CLIENT_ID = "617831923199-ha6054jhlqqkrioohv5fioo5m5f10iki.apps.googleusercontent.com"
const GOOGLE_CLIENT_SECRET = "wtgMIEiAt5UGKjg3BiBNCIf5"
const GOOGLE_REDIRECT_URI = "http://localhost:8080/api/login-google"
const GOOGLE_LINK_GET_TOKEN = "https://accounts.google.com/o/oauth2/token"
const GOOGLE_LINK_GET_USER_INFO = "https://www.googleapis.com/oauth2/v1/userinfo?access_token="
const GOOGLE_GRANT_TYPE = "authorization_code"

// LoginGoogle handler
func LoginGoogle(w http.ResponseWriter, r *http.Request) {
	// Get URL param "code"
	keys, ok := r.URL.Query()["code"]
	if !ok || len(keys[0]) < 1 {
		fmt.Fprintf(w, "URL param 'code' is missing")
	}
	key := keys[0]
	access_token := getToken(key)

	var user_info map[string]interface{}
	user_info = getUserInfo(access_token)

	fmt.Fprintf(w, fmt.Sprintf("%v", map[string]interface{}{
		"id":    user_info["id"],
		"name":  user_info["name"],
		"email": user_info["email"],
	}))
}

func getToken(code string) string {
	data := url.Values{
		"client_id":     {GOOGLE_CLIENT_ID},
		"client_secret": {GOOGLE_CLIENT_SECRET},
		"redirect_uri":  {GOOGLE_REDIRECT_URI},
		"code":          {code},
		"grant_type":    {GOOGLE_GRANT_TYPE},
	}
	resp, err := http.PostForm(GOOGLE_LINK_GET_TOKEN, data)

	if err != nil {
		return ""
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	access_token := fmt.Sprintf("%v", res["access_token"])
	access_token = strings.ReplaceAll(access_token, "\"", "")
	return access_token
}

func getUserInfo(access_token string) map[string]interface{} {
	link := GOOGLE_LINK_GET_USER_INFO + access_token
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
