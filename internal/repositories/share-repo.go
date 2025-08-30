package repositories

import (
	"github.com/asrma7/meroshare-bot/internal/models"
	"gorm.io/gorm"
)

type ShareRepository interface {
	AddAppliedShare(share *models.AppliedShare) error
	AddApplyShareError(error *models.AppliedShareError) error
	GetAppliedSharesByUserID(userID string) ([]models.AppliedShare, error)
	GetAppliedShareErrorsByUserID(userID string) ([]models.AppliedShareError, error)
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

func (s *shareRepository) GetAppliedShareErrorsByUserID(userID string) ([]models.AppliedShareError, error) {
	var errors []models.AppliedShareError
	err := s.db.Where("user_id = ?", userID).Find(&errors).Error
	return errors, err
}
