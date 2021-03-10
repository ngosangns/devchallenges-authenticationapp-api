package api

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	models "github.com/ngosangns/devchallenges-my-unsplash-api/models"
)

// Util handler
func Util(w http.ResponseWriter, r *http.Request) {
	printErr(w, errors.New("404 Not found"), "")
}

func printRes(w http.ResponseWriter, res []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(res)
}

func printErr(w http.ResponseWriter, err error, clientErr string) {
	// Print log
	log.Println(err)
	// Set client message
	if clientErr == "" {
		clientErr = err.Error()
	}
	// Print response
	b, _ := json.Marshal(models.Res{
		Status:  false,
		Message: clientErr,
	})
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}

// createToken return jwt, hash (SHA256)
func createToken(user models.User) (string, string) {
	var header string = b64Encode(`{"alg":"HS256","typ":"JWT"}`)
	var payload string = b64Encode(`{"email":` + user.Email + `}`)
	var signature string = header + "." + payload
	// Assign secret key
	secretHash := sha256.New()
	secretHash.Write([]byte(fmt.Sprintf("%v", time.Now())))
	secretKey := hex.EncodeToString(secretHash.Sum(nil))
	// Assign signature hash string
	jwt := signature + "." + hex.EncodeToString(secretHash.Sum([]byte(signature)))
	return jwt, string(secretKey)
}

func b64Encode(str string) string {
	return b64.StdEncoding.EncodeToString([]byte(str))
}

func regEx(str string, pattern string) bool {
	match, _ := regexp.MatchString(pattern, str)
	return match
}
