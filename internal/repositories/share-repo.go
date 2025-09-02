package repositories

import (
	"github.com/asrma7/meroshare-bot/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShareRepository interface {
	AddAppliedShare(share *models.AppliedShare) (uuid.UUID, error)
	AddApplyShareError(error *models.AppliedShareError) error
	GetAppliedSharesByUserID(userID string) ([]models.AppliedShare, error)
	GetAppliedShareByID(shareID string) (*models.AppliedShare, error)
	GetAppliedShareByAccountIDAndCompanyShareID(accountID string, companyShareID string) (*models.AppliedShare, error)
	GetAppliedShareErrorsByUserID(userID string) ([]models.AppliedShareError, error)
	GetAppliedShareErrorsByAppliedShareID(appliedShareID string) (*models.AppliedShareError, error)
	MarkShareErrorsAsSeenByUserID(userID string) error
	DeleteAllAppliedSharesByUserID(userID uuid.UUID) error
	DeleteAllAppliedShareErrorsByUserID(userID uuid.UUID) error
}

type shareRepository struct {
	db *gorm.DB
}

func NewShareRepository(db *gorm.DB) ShareRepository {
	return &shareRepository{db: db}
}

func (s *shareRepository) AddAppliedShare(share *models.AppliedShare) (uuid.UUID, error) {
	if err := s.db.Create(share).Error; err != nil {
		return uuid.Nil, err
	}
	return share.ID, nil
}

func (s *shareRepository) AddApplyShareError(error *models.AppliedShareError) error {
	return s.db.Create(error).Error
}

func (s *shareRepository) GetAppliedSharesByUserID(userID string) ([]models.AppliedShare, error) {
	var shares []models.AppliedShare
	err := s.db.Where("user_id = ?", userID).Find(&shares).Error
	return shares, err
}

func (s *shareRepository) GetAppliedShareByID(shareID string) (*models.AppliedShare, error) {
	var share models.AppliedShare
	err := s.db.Where("id = ?", shareID).First(&share).Error
	if err != nil {
		return nil, err
	}
	return &share, nil
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

func (s *shareRepository) GetAppliedShareErrorsByAppliedShareID(appliedShareID string) (*models.AppliedShareError, error) {
	var error models.AppliedShareError
	err := s.db.Where("applied_share_id = ?", appliedShareID).First(&error).Error
	if err != nil {
		return nil, err
	}
	return &error, nil
}

func (s *shareRepository) MarkShareErrorsAsSeenByUserID(userID string) error {
	return s.db.Model(&models.AppliedShareError{}).Where("user_id = ?", userID).Update("seen", true).Error
}

func (s *shareRepository) DeleteAllAppliedSharesByUserID(userID uuid.UUID) error {
	return s.db.Where("user_id = ?", userID).Delete(&models.AppliedShare{}).Error
}

func (s *shareRepository) DeleteAllAppliedShareErrorsByUserID(userID uuid.UUID) error {
	return s.db.Where("user_id = ?", userID).Delete(&models.AppliedShareError{}).Error
}
