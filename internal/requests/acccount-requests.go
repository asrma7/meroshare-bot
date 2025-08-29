package requests

type AccountRequest struct {
	ClientId       uint16 `json:"client_id" binding:"required"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	BankId         string `json:"bank_id" binding:"required"`
	CRNNumber      string `json:"crn_number" binding:"required"`
	TransactionPIN string `json:"transaction_pin" binding:"required"`
	PreferredKitta uint16 `json:"preferred_kitta" binding:"required,min=10"`
}
