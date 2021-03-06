package api

import (
	"fmt"
	"net/http"
	"time"
)

// Date handler
func Date(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format(time.RFC850)
	fmt.Fprintf(w, currentTime)
}
