package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type pwnedHandler struct {
	log       *log.Logger
	validator ValidatePassword
}

func NewPwnedHandler(log *log.Logger, validator ValidatePassword) *pwnedHandler {

	return &pwnedHandler{log, validator}
}

func (pw *pwnedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// Retrieve user password from context
	password, err := getPassword(r)
	if err != nil {
		http.Error(rw, "could not retrieve user password", http.StatusInternalServerError)
		return
	}

	// Call to Pwned Service
	isSecure, err := pw.validator(password)
	if err != nil {
		http.Error(rw, "could not verify if password is compromised", http.StatusInternalServerError)
		return
	}

	if !isSecure {
		http.Error(rw, "insecure password", http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func getPassword(r *http.Request) (string, error) {

	password, ok := r.Context().Value(PwnedContextKey("UserPassword")).(string)
	if !ok {
		return "", fmt.Errorf("user password not set in response")
	}
	if (len([]rune(password))) == 0 {
		return "", fmt.Errorf("user password is empty")
	}

	return password, nil
}
