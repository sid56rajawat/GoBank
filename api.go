package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type APIServer struct {
	listenAddress string
	store         Storage
}

func NewAPIServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))
	mux.HandleFunc("GET /account/", makeHTTPHandleFunc(s.handleGetAccount))
	mux.HandleFunc("POST /account", makeHTTPHandleFunc(s.handleCreateAccount))
	mux.HandleFunc("POST /deposit", makeHTTPHandleFunc(s.handleDeposit))
	mux.HandleFunc("POST /withdraw", makeHTTPHandleFunc(s.handleWithdraw))
	mux.HandleFunc("POST /transfer", makeHTTPHandleFunc(s.handleTransfer))
	mux.HandleFunc("DELETE /account/{id}", makeHTTPHandleFunc(s.handleDeleteAccount))

	log.Println("JSON API server running on: ", s.listenAddress)

	http.ListenAndServe(s.listenAddress, mux)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
		return err
	}

	// TODO: bugfix - the account object returned has the wrong id, which doesn't represent the actual id stored in postgres
	account := NewAccount(createAccountReq.FistName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleDeposit(w http.ResponseWriter, r *http.Request) error {
	depositRequest := DepositRequest{}
	if err := json.NewDecoder(r.Body).Decode(&depositRequest); err != nil {
		return err
	}
	if err := s.store.DepositToAccount(depositRequest.ToAccount, depositRequest.Amount); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, "deposit to account successful")
}

func (s *APIServer) handleWithdraw(w http.ResponseWriter, r *http.Request) error {
	withdrawRequest := WithdrawRequest{}
	if err := json.NewDecoder(r.Body).Decode(&withdrawRequest); err != nil {
		return err
	}
	if err := s.store.WithdrawFromAccount(withdrawRequest.FromAccount, withdrawRequest.Amount); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, "withdrawal from account successful")
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferRequest := TransferRequest{}
	if err := json.NewDecoder(r.Body).Decode(&transferRequest); err != nil {
		return err
	}
	if err := s.store.TransferMoney(transferRequest.FromAccount, transferRequest.ToAccount, transferRequest.Amount); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, "money transferred successfully")
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
