package repositories

import (
	"github.com/asrma7/meroshare-bot/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(account *models.Account) (uuid.UUID, error)
	GetAccountByID(id uuid.UUID) (*models.Account, error)
	GetAccountsByUserID(userID uuid.UUID) ([]models.Account, error)
	GetAllAccounts() ([]models.Account, error)
	UpdateAccount(account *models.Account) error
	DeleteAccount(id uuid.UUID) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) CreateAccount(account *models.Account) (uuid.UUID, error) {
	if err := r.db.Create(account).Error; err != nil {
		return uuid.Nil, err
	}
	return account.ID, nil
}

func (r *accountRepository) GetAccountByID(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	if err := r.db.First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetAccountsByUserID(userID uuid.UUID) ([]models.Account, error) {
	var accounts []models.Account
	if err := r.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepository) GetAllAccounts() ([]models.Account, error) {
	var accounts []models.Account
	if err := r.db.Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepository) UpdateAccount(account *models.Account) error {
	if err := r.db.Save(account).Error; err != nil {
		return err
	}
	return nil
}

func (r *accountRepository) DeleteAccount(id uuid.UUID) error {
	if err := r.db.Delete(&models.Account{}, id).Error; err != nil {
		return err
	}
	return nil
}
