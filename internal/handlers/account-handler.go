package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asrma7/meroshare-bot/internal/models"
	"github.com/asrma7/meroshare-bot/internal/requests"
	"github.com/asrma7/meroshare-bot/internal/responses"
	"github.com/asrma7/meroshare-bot/internal/services"
	"github.com/asrma7/meroshare-bot/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type AccountHandler interface {
	CreateAccount(c *gin.Context)
	GetAccountByID(c *gin.Context)
	GetAccountsByUserID(c *gin.Context)
	UpdateAccount(c *gin.Context)
	DeleteAccount(c *gin.Context)
}

type accountHandler struct {
	accountService services.AccountService
}

func NewAccountHandler(accountService services.AccountService) AccountHandler {
	return &accountHandler{
		accountService: accountService,
	}
}

func fetchUserDetails(authorization string) (responses.UserDetails, error) {
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

func fetchBankDetails(authorization string, bankId string) ([]responses.BankDetails, error) {
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

func (h *accountHandler) CreateAccount(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		errResp := errors.ErrorResponse{
			Type:    "UNAUTHORIZED",
			Message: "User ID not found in context",
		}
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		errResp := errors.ErrorResponse{
			Type:    "VALIDATION_ERROR",
			Message: "Invalid user ID format",
		}
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	var req requests.AccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reqData := map[string]string{
		"clientId": fmt.Sprintf("%d", req.ClientId),
		"username": req.Username,
		"password": req.Password,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := http.Post("https://webbackend.cdsc.com.np/api/meroShare/auth/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	authorization := resp.Header.Get("Authorization")
	if authorization == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login account"})
		return
	}

	var (
		userDetails responses.UserDetails
		bankDetails []responses.BankDetails
	)
	g := new(errgroup.Group)

	g.Go(func() error {
		var err error
		userDetails, err = fetchUserDetails(authorization)
		return err
	})

	g.Go(func() error {
		var err error
		bankDetails, err = fetchBankDetails(authorization, fmt.Sprintf("%v", req.BankId))
		return err
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	account := models.Account{
		UserID:             userIDParsed,
		Name:               userDetails.Name,
		Email:              userDetails.Email,
		Contact:            userDetails.Contact,
		ClientID:           req.ClientId,
		Username:           req.Username,
		Password:           req.Password,
		BankID:             req.BankId,
		CRNNumber:          req.CRNNumber,
		TransactionPIN:     req.TransactionPIN,
		AccountTypeId:      bankDetails[0].AccountTypeId,
		PreferredKitta:     fmt.Sprintf("%d", req.PreferredKitta),
		Demat:              userDetails.Demat,
		BOID:               userDetails.BOID,
		AccountNumber:      bankDetails[0].AccountNumber,
		CustomerId:         bankDetails[0].ID,
		AccountBranchId:    bankDetails[0].AccountBranchId,
		DMATExpiryDate:     userDetails.DematExpiryDate,
		PasswordExpiryDate: userDetails.PasswordExpiryDate,
		ExpiredDate:        userDetails.ExpiredDate,
	}

	account_id, err := h.accountService.CreateAccount(&account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Account created successfully", "account_id": account_id})
}

func (h *accountHandler) GetAccountByID(c *gin.Context) {
	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
		return
	}

	parsedAccountID, err := uuid.Parse(accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Account ID"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	account, err := h.accountService.GetAccountByID(parsedAccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account.UserID != userIDParsed {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "account": account})
}

func (h *accountHandler) GetAccountsByUserID(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	accounts, err := h.accountService.GetAccountsByUserID(userIDParsed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "accounts": accounts})
}

func (h *accountHandler) UpdateAccount(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		errResp := errors.ErrorResponse{
			Type:    "UNAUTHORIZED",
			Message: "User ID not found in context",
		}
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		errResp := errors.ErrorResponse{
			Type:    "VALIDATION_ERROR",
			Message: "Invalid user ID format",
		}
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	accountId := c.Param("id")
	if accountId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
		return
	}

	parsedAccountId, err := uuid.Parse(accountId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Account ID"})
		return
	}

	account, err := h.accountService.GetAccountByID(parsedAccountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if account.UserID != userIDParsed {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req requests.AccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reqData := map[string]string{
		"clientId": fmt.Sprintf("%d", req.ClientId),
		"username": req.Username,
		"password": req.Password,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := http.Post("https://webbackend.cdsc.com.np/api/meroShare/auth/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	authorization := resp.Header.Get("Authorization")
	if authorization == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login account"})
		return
	}

	var (
		userDetails responses.UserDetails
		bankDetails []responses.BankDetails
	)
	g := new(errgroup.Group)

	g.Go(func() error {
		var err error
		userDetails, err = fetchUserDetails(authorization)
		return err
	})

	g.Go(func() error {
		var err error
		bankDetails, err = fetchBankDetails(authorization, fmt.Sprintf("%v", req.BankId))
		return err
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedAccount := models.Account{
		ID:                 parsedAccountId,
		UserID:             userIDParsed,
		Name:               userDetails.Name,
		Email:              userDetails.Email,
		Contact:            userDetails.Contact,
		ClientID:           req.ClientId,
		Username:           req.Username,
		Password:           req.Password,
		BankID:             req.BankId,
		CRNNumber:          req.CRNNumber,
		TransactionPIN:     req.TransactionPIN,
		AccountTypeId:      bankDetails[0].AccountTypeId,
		PreferredKitta:     fmt.Sprintf("%d", req.PreferredKitta),
		Demat:              userDetails.Demat,
		BOID:               userDetails.BOID,
		AccountNumber:      bankDetails[0].AccountNumber,
		CustomerId:         bankDetails[0].ID,
		AccountBranchId:    bankDetails[0].AccountBranchId,
		DMATExpiryDate:     userDetails.DematExpiryDate,
		PasswordExpiryDate: userDetails.PasswordExpiryDate,
		ExpiredDate:        userDetails.ExpiredDate,
	}

	err = h.accountService.UpdateAccount(&updatedAccount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Account updated successfully"})
}

func (h *accountHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		errResp := errors.ErrorResponse{
			Type:    "UNAUTHORIZED",
			Message: "User ID not found in context",
		}
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		errResp := errors.ErrorResponse{
			Type:    "VALIDATION_ERROR",
			Message: "Invalid user ID format",
		}
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	accountId := c.Param("id")
	if accountId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
		return
	}

	parsedAccountId, err := uuid.Parse(accountId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Account ID"})
		return
	}

	account, err := h.accountService.GetAccountByID(parsedAccountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account.UserID != userIDParsed {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this account"})
		return
	}

	err = h.accountService.DeleteAccount(parsedAccountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Account deleted successfully"})
}
