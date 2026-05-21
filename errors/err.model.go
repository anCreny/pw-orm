package errors

// INFO: https://learn.microsoft.com/en-us/dotnet/api/system.management.automation.errorrecord?view=powershellsdk-7.4.0
type Error struct {
	CategoryInfo CategoryInfo  `json:"CategoryInfo"`
	ErrorDetails *ErrorDetails `json:"ErrorDetails"`
	Exception    Exception     `json:"Exception"`
}

// never null
// INFO: https://learn.microsoft.com/en-us/dotnet/api/system.management.automation.errorcategoryinfo?view=powershellsdk-7.4.0
type CategoryInfo struct {
	Activity   string `json:"Activity"`
	Category   int    `json:"Category"` // enum
	Reason     string `json:"Reason"`
	TargetName string `json:"TargetName"`
	TargetType string `json:"TargetType"`
}

// may be null
// INFO: https://learn.microsoft.com/en-us/dotnet/api/system.management.automation.errordetails?view=powershellsdk-7.4.0
type ErrorDetails struct {
	Message           string `json:"Message"`
	RecommendedAction string `json:"RecommendedAction"`
}

// never null
// INFO: https://learn.microsoft.com/en-us/dotnet/api/system.exception?view=netframework-4.8.1#properties
type Exception struct {
	Data           any        `json:"Data"`
	HelpLink       *string    `json:"HelpLink"`
	HResult        int        `json:"HResult"`
	InnerException *Exception `json:"InnerException"`
	Message        string     `json:"Message"`
	Source         string     `json:"Source"`
	StackTrace     *string    `json:"StackTrace"`
	TargetSite     any        `json:"TargetSite"`
}
