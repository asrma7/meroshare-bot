package responses

type UserDashboard struct {
	User               UserSummary         `json:"user"`
	TotalAccounts      int                 `json:"total_accounts"`
	AccountsWithIssue  int                 `json:"accounts_with_issue"`
	TotalShares        int                 `json:"total_shares"`
	FailedShares       int                 `json:"failed_shares"`
	FailedApplications []FailedApplication `json:"failed_applications"`
}

type UserSummary struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type FailedApplication struct {
	AccountID   string `json:"account_id"`
	AccountName string `json:"account_name"`
	ShareID     string `json:"share_id"`
	Scrip       string `json:"scrip"`
	ErrorMsg    string `json:"error_message"`
	CreatedAt   string `json:"created_at"`
}
