package api

import (
	"crypto/sha256"
	"fmt"
	"time"
)

//getToken func return generated token for login and log_time
func getToken(login string, genTime time.Time) string {
	tokenHash := sha256.Sum256([]byte(login + genTime.Format("2006-01-02 15:04:05")))
	token := fmt.Sprintf("%x", tokenHash)

	return token
}
