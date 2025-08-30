package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asrma7/meroshare-bot/internal/models"
	"github.com/asrma7/meroshare-bot/internal/repositories"
	"github.com/asrma7/meroshare-bot/internal/responses"
	"github.com/google/uuid"
)

type AccountService interface {
	LoginAccount(clientId uint16, username, password string) (string, error)
	FetchUserDetails(authorization string) (responses.UserDetails, error)
	FetchBankDetails(authorization string, bankId string) ([]responses.BankDetails, error)
	CreateAccount(account *models.Account) (uuid.UUID, error)
	GetAccountByID(id uuid.UUID) (*models.Account, error)
	GetAccountsByUserID(userID uuid.UUID) ([]models.Account, error)
	GetAllAccounts() ([]models.Account, error)
	UpdateAccount(account *models.Account) error
	DeleteAccount(id uuid.UUID) error
}

type accountService struct {
	repo repositories.AccountRepository
}

func NewAccountService(repo *repositories.AccountRepository) AccountService {
	return &accountService{repo: *repo}
}

func (s *accountService) LoginAccount(clientId uint16, username, password string) (string, error) {

	reqData := map[string]string{
		"clientId": fmt.Sprintf("%d", clientId),
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://webbackend.cdsc.com.np/api/meroShare/auth/", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to login account")
	}

	authorization := resp.Header.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("failed to login account")
	}

	return authorization, nil
}

func (s *accountService) FetchUserDetails(authorization string) (responses.UserDetails, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://webbackend.cdsc.com.np/api/meroShare/ownDetail/", nil)
	if err != nil {
		return responses.UserDetails{}, err
	}
	req.Header.Set("Authorization", authorization)

	resp, err := client.Do(req)
	if err != nil {
		return responses.UserDetails{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return responses.UserDetails{}, fmt.Errorf("failed to fetch user details")
	}

	var userDetails responses.UserDetails
	if err := json.NewDecoder(resp.Body).Decode(&userDetails); err != nil {
		return responses.UserDetails{}, err
	}

	return userDetails, nil
}

func (s *accountService) FetchBankDetails(authorization string, bankId string) ([]responses.BankDetails, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://webbackend.cdsc.com.np/api/meroShare/bank/%s", bankId), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authorization)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch bank details")
	}

	var bankDetails []responses.BankDetails
	if err := json.NewDecoder(resp.Body).Decode(&bankDetails); err != nil {
		return nil, err
	}

	return bankDetails, nil
}

func (s *accountService) CreateAccount(account *models.Account) (uuid.UUID, error) {
	return s.repo.CreateAccount(account)
}

func (s *accountService) GetAccountByID(id uuid.UUID) (*models.Account, error) {
	return s.repo.GetAccountByID(id)
}

func (s *accountService) GetAccountsByUserID(userID uuid.UUID) ([]models.Account, error) {
	return s.repo.GetAccountsByUserID(userID)
}

func (s *accountService) GetAllAccounts() ([]models.Account, error) {
	return s.repo.GetAllAccounts()
}

func (s *accountService) UpdateAccount(account *models.Account) error {
	return s.repo.UpdateAccount(account)
}

func (s *accountService) DeleteAccount(id uuid.UUID) error {
	return s.repo.DeleteAccount(id)
}
