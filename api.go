package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	/*We are using gorilla package because in order to handle the request
	  we would need to write a bunch of regular expression to map the request
	  to right handle. This package is very solid an we dont need to code
	  this boiled plate code averytime.
	  So for this part it's ok to use this 3rd party package
	*/
	router := mux.NewRouter()
	router.HandleFunc("/account/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAUTH(makeHTTPHandleFunc(s.handleAccountByID), s.store))
	router.HandleFunc("/account/transfer", makeHTTPHandleFunc(s.handleTransfer))
	fmt.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleAccountByID(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetAccountByID(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {

	loginRequest := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
		return fmt.Errorf("account not found")
	}

	account, err := s.store.GetAccountByNumber(loginRequest.AccountNumber)
	if err != nil {
		return fmt.Errorf("account not found")
	}

	err = account.ValidatePassword(loginRequest.Password)
	if err != nil {
		return fmt.Errorf("account not found")
	}

	token, err := createJWTToken(account)
	if err != nil {
		w.Header().Set("redirect", "/login")
	}

	w.Header().Set("x-jwt-token", token)
	return WriteJSON(w, http.StatusOK, account)
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
		return WriteJSON(w, http.StatusOK, nil)
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	createAccountRequest := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return err
	}

	account, err := NewAccount(createAccountRequest.FisrtName, createAccountRequest.LastName, createAccountRequest.Password)
	if err != nil {
		return fmt.Errorf("an error occured while creating account")
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	token, err := createJWTToken(account)
	if err != nil {
		return err
	}

	response := LoginResponse{
		AccountNumber: account.Number,
		Token:         token,
	}

	return WriteJSON(w, http.StatusOK, response)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"Deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferRequest := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}

	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferRequest)
}

/*
We promote our more import function to the top of the page
and our less important function goes to the bottom as it is a good practice.
*/

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id %s", idStr)
	}
	return id, nil
}

type ApiError struct {
	Error string `json:"error"`
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "Permission denied"})
}
func WriteJSON(w http.ResponseWriter, statusCode int, content any) error {
	//The Content-Type header need to go first
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(content)
}

func withJWTAUTH(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokeString := r.Header.Get("token")
		token, err := validateJWT(tokeString)

		if err != nil {
			fmt.Println("1", err)
			permissionDenied(w)
			return
		}

		id, err := getID(r)
		if err != nil {
			fmt.Println("2", err)
			permissionDenied(w)
			return
		}

		account, err := s.GetAccountByID(id)
		if err != nil {
			fmt.Println("3", err)
			permissionDenied(w)
			return
		}

		if account.Number != token.Claims.(*UserClaims).AccountNumber {
			fmt.Println("4", err)
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		fmt.Println("passow aqui")
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("secret %l", secret)
		hmacSampleSecret := []byte(secret)
		return hmacSampleSecret, nil
	})

	fmt.Println(token, err)
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		fmt.Println("Account number: %D, ExpireAte: %l", claims.AccountNumber, claims.ExpiresAt)
		return token, nil
	} else {
		fmt.Println(err)
		return nil, err
	}
}

func createJWTToken(account *Account) (string, error) {
	// export JWT_SECRET=secrete123
	clams := UserClaims{AccountNumber: account.Number}
	clams.ExpiresAt = jwt.NewNumericDate(time.Now().Add(2 * time.Minute))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clams)

	secret := os.Getenv("JWT_SECRET")
	hmacSampleSecret := []byte(secret)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)

	fmt.Println(tokenString, err)

	return tokenString, err
}

type Apifunc func(w http.ResponseWriter, r *http.Request) error

func MakeHTTPHandleFunc(f Apifunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
