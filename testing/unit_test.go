package testing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jonatasolmartins/gobank"
)

func TestGetAccount(t *testing.T) {

	req := httptest.NewRequest("GET", "/account", nil)

	rr := httptest.NewRecorder()

	store := &MockStorage{
		Accounts: []*gobank.Account{
			{
				ID:                1,
				FisrtName:         "Karina",
				LastName:          "brito",
				Number:            1234,
				EncryptedPassword: "123456",
				Balance:           1000,
				CreatedAt:         time.Now(),
			},
			{
				ID:                2,
				FisrtName:         "Joe",
				LastName:          "Doe",
				Number:            1234,
				EncryptedPassword: "123456",
				Balance:           1000,
				CreatedAt:         time.Now(),
			},
		},
	}

	server := gobank.NewAPIServer(":800", store)

	handler := http.HandlerFunc(gobank.MakeHTTPHandleFunc(server.HandleAccount))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var accounts []gobank.Account
	err := json.NewDecoder(rr.Body).Decode(&accounts)
	if err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}

	accountList, err := store.GetAccounts()
	if err != nil {
		t.Errorf("failed to get accounts: %v", err)
	}

	if len(accounts) != len(accountList) {
		t.Errorf("handler returned unexpected body: got %v want %v", accounts, accountList)
	}

	for i := range accounts {
		account := accountList[i]
		if accounts[i].ID != account.ID ||
			//accounts[i].CreatedAt != account.CreatedAt || //why the time value is diferent?
			accounts[i].Balance != account.Balance {
			t.Errorf("handler returned unexpected body: got %v want %v", accounts, account)
		}
	}

}

type MockStorage struct {
	Accounts []*gobank.Account
}

func (s *MockStorage) CreateAccount(account *gobank.Account) error {
	s.Accounts = append(s.Accounts, account)
	return nil
}

func (s *MockStorage) UpdateAccount(account *gobank.Account) error {
	for i, a := range s.Accounts {
		if a.ID == account.ID {
			s.Accounts[i] = account
			return nil
		}
	}
	return fmt.Errorf("account not found")
}

func (s *MockStorage) GetAccountByID(id int) (*gobank.Account, error) {
	for _, a := range s.Accounts {
		if a.ID == int64(id) {
			return a, nil
		}
	}
	return nil, fmt.Errorf("account not found")
}

func (s *MockStorage) GetAccounts() ([]*gobank.Account, error) {
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

func (s *MockStorage) GetAccountByNumber(number int64) (*gobank.Account, error) {
	for _, a := range s.Accounts {
		if a.Number == number {
			return a, nil
		}
	}
	return nil, fmt.Errorf("account not found")
}
