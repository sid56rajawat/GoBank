package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddress string
}

func NewAPIServer(listenAddress string) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()

	// mux.HandleFunc("GET /account", makeHTTPHandleFunc(s.handleGetAccount))
	mux.HandleFunc("GET /account/{id}", makeHTTPHandleFunc(s.handleGetAccount))
	mux.HandleFunc("POST /account", makeHTTPHandleFunc(s.handleCreateAccount))
	mux.HandleFunc("DELETE /account", makeHTTPHandleFunc(s.handleDeleteAccount))

	log.Println("JSON API server running on: ", s.listenAddress)

	http.ListenAndServe(s.listenAddress, mux)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	account := NewAccount("Sid", "Rajawat")
	log.Println("Fetching Details of Account No:", id)

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
