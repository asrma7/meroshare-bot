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
	// GetAccountByID(c *gin.Context)
	// GetAccountByUserID(c *gin.Context)
	// GetAllAccounts(c *gin.Context)
	// UpdateAccount(c *gin.Context)
	// DeleteAccount(c *gin.Context)
}

type accountHandler struct {
	accountService services.AccountService
}

func NewAccountHandler(accountService services.AccountService) AccountHandler {
	return &accountHandler{
		accountService: accountService,
	}
}

func (h *accountHandler) CreateAccount(c *gin.Context) {
	userID := c.MustGet("userID").(string)
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
	var req requests.CreateAccountRequest
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

	c.JSON(http.StatusOK, gin.H{"message": "Account created successfully", "account_id": account_id})
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
