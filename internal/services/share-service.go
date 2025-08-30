package services

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/asrma7/meroshare-bot/internal/models"
	"github.com/asrma7/meroshare-bot/internal/repositories"
	"github.com/asrma7/meroshare-bot/internal/requests"
	"github.com/asrma7/meroshare-bot/internal/responses"
)

type ShareService interface {
	AddAppliedShare(share *models.AppliedShare) error
	AddApplyShareError(error *models.AppliedShareError) error
	GetAppliedSharesByUserID(userID string) ([]models.AppliedShare, error)
	GetAppliedShareErrorsByUserID(userID string) ([]models.AppliedShareError, error)
	FetchApplicableShares(authorization string) (responses.ApplicableSharesResponse, error)
	ApplyForShare(account models.Account, share responses.ApplicableShare, authorization string) (map[string]any, error)
}

type shareService struct {
	repo repositories.ShareRepository
}

func NewShareService(repo *repositories.ShareRepository) ShareService {
	return &shareService{repo: *repo}
}

func (s *shareService) AddAppliedShare(share *models.AppliedShare) error {
	return s.repo.AddAppliedShare(share)
}

func (s *shareService) AddApplyShareError(error *models.AppliedShareError) error {
	return s.repo.AddApplyShareError(error)
}

func (s *shareService) GetAppliedSharesByUserID(userID string) ([]models.AppliedShare, error) {
	return s.repo.GetAppliedSharesByUserID(userID)
}

func (s *shareService) GetAppliedShareErrorsByUserID(userID string) ([]models.AppliedShareError, error) {
	return s.repo.GetAppliedShareErrorsByUserID(userID)
}

func (s *shareService) FetchApplicableShares(authorization string) (responses.ApplicableSharesResponse, error) {
	client := http.Client{}
	payload := strings.NewReader(`{
    "filterFieldParams": [
        {
            "key": "companyIssue.companyISIN.script",
            "alias": "Scrip"
        },
        {
            "key": "companyIssue.companyISIN.company.name",
            "alias": "Company Name"
        },
        {
            "key": "companyIssue.assignedToClient.name",
            "value": "",
            "alias": "Issue Manager"
        }
    ],
    "page": 1,
    "size": 10,
    "searchRoleViewConstants": "VIEW_APPLICABLE_SHARE",
    "filterDateParams": [
        {
            "key": "minIssueOpenDate",
            "condition": "",
            "alias": "",
            "value": ""
        },
        {
            "key": "maxIssueCloseDate",
            "condition": "",
            "alias": "",
            "value": ""
        }
    ]
}`)
	req, err := http.NewRequest("POST", "https://webbackend.cdsc.com.np/api/meroShare/companyShare/applicableIssue/", payload)
	if err != nil {
		return responses.ApplicableSharesResponse{}, err
	}
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return responses.ApplicableSharesResponse{}, err
	}
	defer resp.Body.Close()

	var applicableShares responses.ApplicableSharesResponse
	if err := json.NewDecoder(resp.Body).Decode(&applicableShares); err != nil {
		return responses.ApplicableSharesResponse{}, err
	}

	return applicableShares, nil
}

func (s *shareService) ApplyForShare(account models.Account, share responses.ApplicableShare, authorization string) (map[string]any, error) {
	req := requests.ApplyShareRequest{
		Demat:           account.Demat,
		BOID:            account.BOID,
		AccountNumber:   account.AccountNumber,
		CustomerID:      account.CustomerId,
		AccountBranchID: account.AccountBranchId,
		AccountTypeID:   account.AccountTypeId,
		AppliedKitta:    account.PreferredKitta,
		CRNNumber:       account.CRNNumber,
		TransactionPIN:  account.TransactionPIN,
		CompanyShareID:  fmt.Sprintf("%d", share.CompanyShareID),
		BankID:          account.BankID,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Error marshaling request: %v", err)
	}

	client := http.Client{}
	httpReq, err := http.NewRequest("POST", "https://webbackend.cdsc.com.np/api/meroShare/applicantForm/share/apply/", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Authorization", authorization)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusConflict {
			if resp.Header.Get("Content-Type") == "application/xml" || strings.Contains(resp.Header.Get("Content-Type"), "application/xml") {
				var xmlErr struct {
					Message string `xml:"message"`
				}
				if err := xml.NewDecoder(resp.Body).Decode(&xmlErr); err != nil {
					return nil, fmt.Errorf("failed to apply for share: %d", resp.StatusCode)
				}
				if xmlErr.Message == "You have entered wrong transaction PIN." {
					return nil, fmt.Errorf("invalid transaction PIN")
				}
				return nil, fmt.Errorf("conflict: %s", xmlErr.Message)
			}
			if resp.Header.Get("Content-Type") == "application/json" || strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
				var jsonErr struct {
					Message string `json:"message"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&jsonErr); err != nil {
					return nil, fmt.Errorf("failed to apply for share: %d", resp.StatusCode)
				}
				if jsonErr.Message == "Application in process. Please try again later." {
					fmt.Println("Application in process. Skipping duplicate application.")
					return nil, nil
				}
				return nil, fmt.Errorf("conflict: %s", jsonErr.Message)
			}
			return nil, fmt.Errorf("failed to apply for share: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("failed to apply for share: %d", resp.StatusCode)
	}

	var responseBody map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, err
	}

	return responseBody, nil
}
