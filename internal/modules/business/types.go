package business

type UpdateProfileRequest struct {
	BizName   string `json:"bizName"`
	BizType   string `json:"bizType"`
	BizSlug   string `json:"bizSlug"`
	Timezone  string `json:"timezone"`
	Country   string `json:"country"`
}

type UpdateWhatsAppRequest struct {
	OwnerWA string `json:"ownerWa"`
	WANum   string `json:"waNum"`
}

type UpdateSettingsRequest struct {
	AutomationEnabled      *bool  `json:"automationEnabled"`
	DefaultInterval        *int   `json:"defaultInterval"`
	SendTime               string `json:"sendTime"`
	Timezone               string `json:"timezone"`
	BillingNotifLow        *bool  `json:"billingNotifLow"`
	BillingNotifCritical   *bool  `json:"billingNotifCritical"`
	BillingNotifSubLow     *bool  `json:"billingNotifSubLow"`
	BillingNotifPreRenewal *bool  `json:"billingNotifPreRenewal"`
	AutoTopupEnabled       *bool  `json:"autoTopup"`
	AutoTopupThreshold     *int   `json:"autoTopupThreshold"`
	AutoTopupPackageID     string `json:"autoTopupPackage"`
}
