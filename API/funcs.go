package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
)

func UniversalParseForm(w *http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()

	if err != nil {
		(*w).WriteHeader(http.StatusBadRequest)
		(*w).Write([]byte("BAD REQUEST: Parse form error"))
	}

	return err
}

func WhereBlock(valName string, value string, where *string) {
	if value == "NULL" {
		if (*where) != "" {
			(*where) += " AND "
		}
		(*where) += valName + " is NULL"

		return
	}

	if value != "" {
		if (*where) != "" {
			(*where) += " AND "
		}
		(*where) += valName + " in (" + value + ")"
	}
}

func SetBlock(valName string, value string, set *string, isString bool) {
	if value != "" {
		if (*set) != "" {
			(*set) += ", "
		}

		if isString {
			value = "'" + value + "'"
		}
		(*set) += valName + " = " + value
	}
}

func GetToken(login string, genTime time.Time) string {
	tokenHash := sha256.Sum256([]byte(login + genTime.Format("2006-01-02 15:04:05")))
	token := fmt.Sprintf("%x", tokenHash)

	return token
}
