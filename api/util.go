package api

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	models "github.com/ngosangns/devchallenges-my-unsplash-api/models"
)

func printErr(w http.ResponseWriter, err error) {
	// Print log
	log.Println(err)
	// Print response
	b, _ := json.Marshal(models.Res{
		Status:  false,
		Message: err.Error(),
	})
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
