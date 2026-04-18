package customer

type ServiceInput struct {
	CategoryID string `json:"categoryId"`
	Date       string `json:"date"`
}

type CreateCustomerRequest struct {
	Name        string         `json:"name"`
	PhoneNumber string         `json:"phoneNumber"`
	Via         string         `json:"via"`
	Services    []ServiceInput `json:"services"`
}

type UpdateCustomerRequest struct {
	Name        *string        `json:"name"`
	PhoneNumber *string        `json:"phoneNumber"`
	Via         *string        `json:"via"`
	Services    []ServiceInput `json:"services"`
}

type VisitRequest struct {
	CustomerID          string   `json:"customerId"`
	CustomerName        string   `json:"customerName"`
	CustomerPhoneNumber string   `json:"customerPhoneNumber"`
	Date                string   `json:"date"`
	CategoryIDs         []string `json:"categoryIds"`
}

type CheckinLookupRequest struct {
	PhoneNumber string `json:"phoneNumber"`
}

type CheckinSubmitRequest struct {
	PhoneNumber string   `json:"phoneNumber"`
	Name        string   `json:"name"`
	Date        string   `json:"date"`
	CategoryIDs []string `json:"categoryIds"`
}

type PaginationMeta struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"totalPages"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev"`
}

type StatusCount struct {
	Semua     int `json:"semua"`
	Aktif     int `json:"aktif"`
	Mendekati int `json:"mendekati"`
	Hilang    int `json:"hilang"`
}

type ListCustomersResponse struct {
	Data        []map[string]interface{} `json:"data"`
	Pagination  PaginationMeta           `json:"pagination"`
	StatusCount StatusCount              `json:"statusCount"`
}
