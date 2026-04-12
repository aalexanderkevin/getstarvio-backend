package customer

type ServiceInput struct {
	CategoryID string `json:"categoryId"`
	Name       string `json:"name"`
	Date       string `json:"date"`
}

type CreateCustomerRequest struct {
	Name     string         `json:"name"`
	WA       string         `json:"wa"`
	Via      string         `json:"via"`
	Services []ServiceInput `json:"services"`
}

type UpdateCustomerRequest struct {
	Name     *string        `json:"name"`
	WA       *string        `json:"wa"`
	Via      *string        `json:"via"`
	Services []ServiceInput `json:"services"`
}

type VisitRequest struct {
	CustomerID   string   `json:"customerId"`
	CustomerName string   `json:"customerName"`
	CustomerWA   string   `json:"customerWa"`
	Date         string   `json:"date"`
	CategoryIDs  []string `json:"categoryIds"`
}

type CheckinLookupRequest struct {
	WA string `json:"wa"`
}

type CheckinSubmitRequest struct {
	WA          string   `json:"wa"`
	Name        string   `json:"name"`
	Date        string   `json:"date"`
	CategoryIDs []string `json:"categoryIds"`
}
