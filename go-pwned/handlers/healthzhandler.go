package handlers

import (
	"log"
	"net/http"
)

type (
	healthzHandler struct {
		l *log.Logger
	}
)

func NewHealthzHandler(log *log.Logger) *healthzHandler {
	return &healthzHandler{log}
}

func (hh *healthzHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		rw.WriteHeader(http.StatusOK)
	}
}
