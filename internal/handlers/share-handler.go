package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/asrma7/meroshare-bot/internal/models"
	"github.com/asrma7/meroshare-bot/internal/services"
	"github.com/asrma7/meroshare-bot/pkg/errors"
	"github.com/asrma7/meroshare-bot/pkg/logs"
	"github.com/asrma7/meroshare-bot/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ShareHandler interface {
	GetAppliedShares(c *gin.Context)
	GetAppliedShareErrors(c *gin.Context)
	GetAppliedShareByID(c *gin.Context)
	ApplyShare()
}

type shareHandler struct {
	shareService   services.ShareService
	accountService services.AccountService
}

func NewShareHandler(shareService services.ShareService, accountService services.AccountService) ShareHandler {
	return &shareHandler{
		shareService:   shareService,
		accountService: accountService,
	}
}

func (h *shareHandler) GetAppliedShares(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	if userID == "" {
		errResp := errors.ErrorResponse{
			Type:    "UNAUTHORIZED",
			Message: "User ID not found in context",
		}
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}

	appliedShares, err := h.shareService.GetAppliedSharesByUserID(userID)
	if err != nil {
		errorResp, statusCode := errors.GetErrorResponse(err)
		c.JSON(statusCode, errorResp)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "applied_shares": appliedShares})
}

func (h *shareHandler) GetAppliedShareErrors(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	if userID == "" {
		errResp := errors.ErrorResponse{
			Type:    "UNAUTHORIZED",
			Message: "User ID not found in context",
		}
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}

	appliedShareErrors, err := h.shareService.GetAppliedShareErrorsByUserID(userID)
	if err != nil {
		errorResp, statusCode := errors.GetErrorResponse(err)
		c.JSON(statusCode, errorResp)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "applied_share_errors": appliedShareErrors})
}

func (h *shareHandler) GetAppliedShareByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		errResp := errors.ErrorResponse{
			Type:    "BAD_REQUEST",
			Message: "Invalid share ID",
		}
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	appliedShare, appliedShareError, err := h.shareService.GetAppliedShareByID(id)
	if err != nil {
		errorResp, statusCode := errors.GetErrorResponse(err)
		c.JSON(statusCode, errorResp)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "applied_share": appliedShare, "applied_share_error": appliedShareError})
}

func (h *shareHandler) ApplyShare() {
	allAccounts, err := h.accountService.GetAllAccounts()
	if err != nil {
		logs.Error("Failed to get all accounts", map[string]any{"error": err})
		return
	}

	for _, account := range allAccounts {
		if account.Status != "active" {
			continue
		}
		if account.ExpiredDate.Before(time.Now()) {
			h.accountService.SetAccountStatus(account.ID, "meroshare_expired")
			continue
		}
		if account.PasswordExpiryDate.Before(time.Now()) {
			h.accountService.SetAccountStatus(account.ID, "password_expired")
			continue
		}
		bsDate := strings.Split(account.DMATExpiryDate, "-")
		dmatExpiryDate, err := utils.ConvertBSToAD(utils.StringToInt(bsDate[0]), utils.StringToInt(bsDate[1]), utils.StringToInt(bsDate[2]))
		if err != nil {
			logs.Error("Failed to convert DMAT expiry date", map[string]any{"error": err})
			continue
		}
		if dmatExpiryDate.Before(time.Now()) {
			h.accountService.SetAccountStatus(account.ID, "dmat_expired")
			continue
		}
		authorization, err := h.accountService.LoginAccount(account.ClientID, account.Username, account.Password)
		if err != nil {
			if err.Error() == "invalid credentials" {
				h.accountService.SetAccountStatus(account.ID, "invalid_credentials")
				continue
			}
			logs.Error("Failed to get authorization header", map[string]any{"error": err})
			continue
		}
		applicableShares, err := h.shareService.FetchApplicableShares(authorization)
		if err != nil {
			logs.Error("Failed to fetch applicable shares", map[string]any{"error": err})
			continue
		}
		for _, share := range applicableShares.Shares {
			alreadyApplied, err := h.shareService.CheckIfShareAlreadyApplied(account.ID.String(), fmt.Sprintf("%d", share.CompanyShareID))
			if err != nil {
				logs.Error("Failed to check if share already applied", map[string]any{"error": err})
				continue
			}
			if share.Action == "" && share.ShareGroupName == "Ordinary Shares" && !alreadyApplied {
				result, err := h.shareService.ApplyForShare(account, share, authorization)
				if err != nil {
					if err.Error() == "conflict: You have entered wrong transaction PIN." {
						h.accountService.SetAccountStatus(account.ID, "invalid_pin")
					}
					logs.Error("Failed to apply for share", map[string]any{"error": err})
					appliedShareId, er := h.shareService.AddAppliedShare(&models.AppliedShare{
						UserID:         account.UserID,
						AccountID:      account.ID,
						CompanyName:    share.CompanyName,
						CompanyShareID: share.CompanyShareID,
						Scrip:          share.Scrip,
						AppliedKitta:   account.PreferredKitta,
						ShareGroupName: share.ShareGroupName,
						ShareTypeName:  share.ShareTypeName,
						SubGroup:       share.SubGroup,
						Status:         "failed",
					})
					if er != nil {
						logs.Error("Failed to add applied share", map[string]any{"error": er})
					}
					h.shareService.AddApplyShareError(&models.AppliedShareError{
						UserID:         account.UserID,
						AccountID:      account.ID,
						AppliedShareID: appliedShareId,
						Message:        err.Error(),
					})
					continue
				}
				logs.Info("Successfully applied for share", map[string]any{"result": result, "account_id": account.ID, "share_id": share.CompanyShareID})
				h.shareService.AddAppliedShare(&models.AppliedShare{
					UserID:         account.UserID,
					AccountID:      account.ID,
					CompanyName:    share.CompanyName,
					CompanyShareID: share.CompanyShareID,
					Scrip:          share.Scrip,
					AppliedKitta:   account.PreferredKitta,
					ShareGroupName: share.ShareGroupName,
					ShareTypeName:  share.ShareTypeName,
					SubGroup:       share.SubGroup,
					Status:         "applied",
				})
			}
		}
	}
}
