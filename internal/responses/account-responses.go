package responses

import "time"

type BankDetails struct {
	AccountBranchId uint32 `json:"accountBranchId"`
	AccountNumber   string `json:"accountNumber"`
	AccountTypeId   uint8  `json:"accountTypeId"`
	AccountTypeName string `json:"accountTypeName"`
	BranchName      string `json:"branchName"`
	ID              uint32 `json:"id"`
}

type UserDetails struct {
	BOID               string    `json:"boid"`
	Contact            string    `json:"contact"`
	Demat              string    `json:"demat"`
	Email              string    `json:"email"`
	Name               string    `json:"name"`
	DematExpiryDate    string    `json:"dematExpiryDate"`
	PasswordExpiryDate time.Time `json:"passwordExpiryDate"`
	ExpiredDate        time.Time `json:"expiredDate"`
}
