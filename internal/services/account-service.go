package services

import (
	"github.com/asrma7/meroshare-bot/internal/models"
	"github.com/asrma7/meroshare-bot/internal/repositories"
	"github.com/google/uuid"
)

type AccountService interface {
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
