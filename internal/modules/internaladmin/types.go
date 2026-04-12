package internaladmin

type PlanConfigResponse struct {
	BusinessID string `json:"businessId"`
	FreeBonus  int    `json:"freeBonus"`
	SubCredits int    `json:"subCredits"`
	SubPrice   int    `json:"subPrice"`
	TopupPrice int    `json:"topupPrice"`
	Tier1Price int    `json:"tier1Price"`
	Tier1Creds int    `json:"tier1Credits"`
	Tier2Price int    `json:"tier2Price"`
	Tier2Creds int    `json:"tier2Credits"`
	Tier3Price int    `json:"tier3Price"`
	Tier3Creds int    `json:"tier3Credits"`
}
