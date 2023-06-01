package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser(t *testing.T) {
	// Create a new request with a dummy user ID
	req := httptest.NewRequest("GET", "/account", nil)

	// Create a new response recorder to capture the response
	rr := httptest.NewRecorder()

	store := &MockStorage{}

	server := NewAPIServer(":800", store)
	// Call the GetUser function with the request and response recorder
	handler := http.HandlerFunc(MakeHTTPHandleFunc(server.HandleAccount))
	handler.ServeHTTP(rr, req)

	// Check that the status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check that the response body contains the expected user data
	var accounts []Account
	err := json.NewDecoder(rr.Body).Decode(&accounts)
	if err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}

	accountList, err := store.GetAccounts()
	if err != nil {
		t.Errorf("failed to get accounts: %v", err)
	}

	if !compareAccounts(accounts, accountList) {
		t.Errorf("handler returned unexpected user data: got %v want %v",
			accounts, accountList)
	}
}

func compareAccounts(a []Account, b &[]Account) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].ID != b[i].ID ||
			a[i].CreatedAt != b[i].CreatedAt ||
			a[i].Balance != b[i].Balance {
			return false
		}
	}

	return true
}

type MockStorage struct {
	Accounts []*Account
}

func (s *MockStorage) CreateAccount(account *Account) error {
	s.Accounts = append(s.Accounts, account)
	return nil
}

func (s *MockStorage) UpdateAccount(account *Account) error {
	for i, a := range s.Accounts {
		if a.ID == account.ID {
			s.Accounts[i] = account
			return nil
		}
	}
	return fmt.Errorf("account not found")
}

func (s *MockStorage) GetAccountByID(id int) (*Account, error) {
	for _, a := range s.Accounts {
		if a.ID == int64(id) {
			return a, nil
		}
	}
	return nil, fmt.Errorf("account not found")
}

func (s *MockStorage) GetAccounts() ([]*Account, error) {
	return s.Accounts, nil
}

func (s *MockStorage) DeleteAccount(id int) error {
	for i, a := range s.Accounts {
		if a.ID == int64(id) {
			s.Accounts = append(s.Accounts[:i], s.Accounts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("account not found")
}

func (s *MockStorage) GetAccountByNumber(number int64) (*Account, error) {
	for _, a := range s.Accounts {
		if a.Number == number {
			return a, nil
		}
	}
	return nil, fmt.Errorf("account not found")
}
