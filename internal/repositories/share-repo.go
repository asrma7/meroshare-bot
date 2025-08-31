package repositories

import (
	"github.com/asrma7/meroshare-bot/internal/models"
	"gorm.io/gorm"
)

type ShareRepository interface {
	AddAppliedShare(share *models.AppliedShare) error
	AddApplyShareError(error *models.AppliedShareError) error
	GetAppliedSharesByUserID(userID string) ([]models.AppliedShare, error)
	GetAppliedShareByAccountIDAndCompanyShareID(accountID string, companyShareID string) (*models.AppliedShare, error)
	GetAppliedShareErrorsByUserID(userID string) ([]models.AppliedShareError, error)
	GetUnseenShareErrorsByUserID(userID string) ([]models.AppliedShareError, error)
	MarkShareErrorsAsSeenByUserID(userID string) error
}

type shareRepository struct {
	db *gorm.DB
}

func NewShareRepository(db *gorm.DB) ShareRepository {
	return &shareRepository{db: db}
}

func (s *shareRepository) AddAppliedShare(share *models.AppliedShare) error {
	return s.db.Create(share).Error
}

func (s *shareRepository) AddApplyShareError(error *models.AppliedShareError) error {
	return s.db.Create(error).Error
}

func (s *shareRepository) GetAppliedSharesByUserID(userID string) ([]models.AppliedShare, error) {
	var shares []models.AppliedShare
	err := s.db.Where("user_id = ?", userID).Find(&shares).Error
	return shares, err
}

func (s *shareRepository) GetAppliedShareByAccountIDAndCompanyShareID(accountID string, companyShareID string) (*models.AppliedShare, error) {
	var share models.AppliedShare
	err := s.db.Where("account_id = ? AND company_share_id = ?", accountID, companyShareID).First(&share).Error
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (s *shareRepository) GetAppliedShareErrorsByUserID(userID string) ([]models.AppliedShareError, error) {
	var errors []models.AppliedShareError
	err := s.db.Where("user_id = ?", userID).Find(&errors).Error
	return errors, err
}

func (s *shareRepository) GetUnseenShareErrorsByUserID(userID string) ([]models.AppliedShareError, error) {
	var errors []models.AppliedShareError
	err := s.db.Where("user_id = ? AND seen = ?", userID, false).Find(&errors).Error
	return errors, err
}

func (s *shareRepository) MarkShareErrorsAsSeenByUserID(userID string) error {
	return s.db.Model(&models.AppliedShareError{}).Where("user_id = ?", userID).Update("seen", true).Error
}
