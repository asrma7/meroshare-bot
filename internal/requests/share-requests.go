package requests

// {
//     "demat": "1301580001737395",
//     "boid": "01737395",
//     "accountNumber": "023011060008386",
//     "customerId": 6050709,
//     "accountBranchId": 4575,
//     "accountTypeId": 1,
//     "appliedKitta": "10",
//     "crnNumber": "SNMAR001030199",
//     "transactionPIN": "1997",
//     "companyShareId": "710",
//     "bankId": "54"
// }

type ApplyShareRequest struct {
	Demat           string `json:"demat"`
	BOID            string `json:"boid"`
	AccountNumber   string `json:"accountNumber"`
	CustomerID      uint32 `json:"customerId"`
	AccountBranchID uint32 `json:"accountBranchId"`
	AccountTypeID   uint8  `json:"accountTypeId"`
	AppliedKitta    string `json:"appliedKitta"`
	CRNNumber       string `json:"crnNumber"`
	TransactionPIN  string `json:"transactionPIN"`
	CompanyShareID  string `json:"companyShareId"`
	BankID          string `json:"bankId"`
}
