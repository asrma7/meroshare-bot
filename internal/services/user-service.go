package services

import (
	"github.com/asrma7/meroshare-bot/internal/models"
	"github.com/asrma7/meroshare-bot/internal/responses"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	GetUserDashboard(userID uuid.UUID) (responses.UserDashboard, error)
	ResetUserLogs(userID uuid.UUID) error
}

type userService struct {
	db           *gorm.DB
	shareService ShareService
}

func NewUserService(db *gorm.DB, shareService ShareService) UserService {
	return &userService{
		db:           db,
		shareService: shareService,
	}
}

func (s *userService) GetUserDashboard(userID uuid.UUID) (responses.UserDashboard, error) {
	var resp responses.UserDashboard
	var user models.User
	if err := s.db.Select("first_name", "last_name").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		return responses.UserDashboard{}, err
	}
	resp.User = responses.UserSummary{
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	var totalAccounts int64
	if err := s.db.Model(&models.Account{}).
		Where("user_id = ?", userID).
		Count(&totalAccounts).Error; err != nil {
		return responses.UserDashboard{}, err
	}
	resp.TotalAccounts = int(totalAccounts)

	var accountsWithIssue int64
	if err := s.db.Model(&models.Account{}).
		Where("user_id = ? AND status <> ?", userID, "active").
		Count(&accountsWithIssue).Error; err != nil {
		return responses.UserDashboard{}, err
	}
	resp.AccountsWithIssue = int(accountsWithIssue)

	var totalShares int64
	if err := s.db.Model(&models.AppliedShare{}).
		Where("user_id = ?", userID).
		Count(&totalShares).Error; err != nil {
		return responses.UserDashboard{}, err
	}
	resp.TotalShares = int(totalShares)

	var failedShares int64
	if err := s.db.Model(&models.AppliedShare{}).
		Where("user_id = ? AND status = ?", userID, "failed").
		Count(&failedShares).Error; err != nil {
		return responses.UserDashboard{}, err
	}
	resp.FailedShares = int(failedShares)

	type failedRow struct {
		AccountID   uuid.UUID
		AccountName string
		ShareID     uuid.UUID
		Scrip       string
		ErrorMsg    string
		CreatedAt   string
	}

	var rows []failedRow
	if err := s.db.Table("applied_share_errors ase").
		Select(`
			ase.account_id,
			a.name AS account_name,
			ase.applied_share_id AS share_id,
			as_.scrip,
			ase.message AS error_msg,
			ase.created_at
		`).
		Joins("JOIN accounts a ON ase.account_id = a.id").
		Joins("JOIN applied_shares as_ ON ase.applied_share_id = as_.id").
		Where("ase.user_id = ?", userID).
		Where("ase.seen = ?", false).
		Order("ase.created_at DESC").
		Scan(&rows).Error; err != nil {
		return responses.UserDashboard{}, err
	}

	for _, r := range rows {
		resp.FailedApplications = append(resp.FailedApplications, responses.FailedApplication{
			AccountID:   r.AccountID.String(),
			AccountName: r.AccountName,
			ShareID:     r.ShareID.String(),
			Scrip:       r.Scrip,
			ErrorMsg:    r.ErrorMsg,
			CreatedAt:   r.CreatedAt,
		})
	}

	return resp, nil
}

func (s *userService) ResetUserLogs(userID uuid.UUID) error {
	err := s.shareService.DeleteAllAppliedShareErrorsByUserID(userID)
	if err != nil {
		return err
	}
	err = s.shareService.DeleteAllAppliedSharesByUserID(userID)
	if err != nil {
		return err
	}
	return nil
}
