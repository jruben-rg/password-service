package handlers

import (
	"context"
	"encoding/base64"
	json "encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	JobName = "PasswordValidations"
)

type (
	PwnedContextKey string

	ValidatePassword func(password string) (bool, error)

	passwordHandler struct {
		l                *log.Logger
		validatePassword ValidatePassword
		next             http.Handler
	}

	passwordRequest struct {
		Password string `json:"password"`
	}
)

func NewPasswordHandler(l *log.Logger, validator ValidatePassword, handler http.Handler) *passwordHandler {
	return &passwordHandler{l, validator, handler}
}

func (ph *passwordHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		start := time.Now()
		decoder := json.NewDecoder(r.Body)
		passwordRequest := passwordRequest{}
		err := decoder.Decode(&passwordRequest)
		if err != nil {
			http.Error(rw, "error decoding request", http.StatusBadRequest)
			return
		}

		decodedPassword, err := decodePassword(passwordRequest.Password)
		if err != nil {
			http.Error(rw, "could not decode password value", http.StatusBadRequest)
			return
		}

		ok, err := ph.validatePassword(decodedPassword)
		if !ok {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			elapsed := time.Since(start)
			log.Printf("Request took %s", elapsed)
			return
		}

		elapsed := time.Since(start)
		log.Printf("PasswordHandler took %s", elapsed)

		//if at this stage all validators are correct, invoke next handler
		if ph.next != nil {
			context := context.WithValue(r.Context(), PwnedContextKey("UserPassword"), decodedPassword)
			r = r.WithContext(context)
			ph.next.ServeHTTP(rw, r)
		}

	}
}

func decodePassword(encodedPassword string) (string, error) {
	decodedValue, err := base64.StdEncoding.DecodeString(encodedPassword)
	if err != nil {
		return "", fmt.Errorf("error decoding password value %s", err)
	}
	return string(decodedValue), nil
}
