package responses

type ApplicableShare struct {
	CompanyShareID uint16 `json:"companyShareId"`
	SubGroup       string `json:"subGroup"`
	Scrip          string `json:"scrip"`
	CompanyName    string `json:"companyName"`
	ShareTypeName  string `json:"shareTypeName"`
	ShareGroupName string `json:"shareGroupName"`
	StatusName     string `json:"statusName"`
	Action         string `json:"action"`
	IssueOpenDate  string `json:"issueOpenDate"`
	IssueCloseDate string `json:"issueCloseDate"`
}

type ApplicableSharesResponse struct {
	Shares     []ApplicableShare `json:"object"`
	TotalCount int               `json:"totalCount"`
}
