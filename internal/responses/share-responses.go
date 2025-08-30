package responses

//	{
//	    "object": [
//	        {
//	            "companyShareId": 710,
//	            "subGroup": "For General Public",
//	            "scrip": "BCTL",
//	            "companyName": "Bandipur Cable Car and Tourism Ltd",
//	            "shareTypeName": "IPO",
//	            "shareGroupName": "Ordinary Shares",
//	            "statusName": "CREATE_APPROVE",
//	            "action": "inProcess",
//	            "issueOpenDate": "Aug 27, 2025 10:00:00 AM",
//	            "issueCloseDate": "Aug 31, 2025 5:00:00 PM"
//	        }
//	    ],
//	    "totalCount": 0
//	}
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
