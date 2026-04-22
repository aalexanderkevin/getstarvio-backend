package routes

// ErrorResponseDoc describes standard error response payload.
type ErrorResponseDoc struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"invalid request payload"`
}

// OKDataDoc is common success marker.
type OKDataDoc struct {
	OK bool `json:"ok" example:"true"`
}

// OKResponseDoc describes standard success response for mutation endpoints.
type OKResponseDoc struct {
	Error bool      `json:"error" example:"false"`
	Data  OKDataDoc `json:"data"`
}

// HealthResponse is the health check payload.
type HealthResponse struct {
	OK      bool   `json:"ok" example:"true"`
	Service string `json:"service" example:"getstarvio-backend"`
	Time    string `json:"time" example:"2026-04-12T10:00:00Z"`
}

type AuthGoogleLoginRequestDoc struct {
	IDToken   string `json:"idToken" example:"eyJhbGciOiJSUzI1NiIs..."`
	GoogleSub string `json:"googleSub" example:"109876543210987654321"`
	Email     string `json:"email" example:"owner@getstarvio.com"`
	Name      string `json:"name" example:"Alexandra Kevin"`
}

type AuthRefreshRequestDoc struct {
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIs..."`
}

type AuthLogoutRequestDoc struct {
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIs..."`
}

type AuthTokenBodyDoc struct {
	AccessToken  string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIs..."`
	UserID       string `json:"userId" example:"7f810c87-1528-4f67-b6c8-0ca01d08860c"`
}

type AuthTokenResponseDoc struct {
	Error bool             `json:"error" example:"false"`
	Data  AuthTokenBodyDoc `json:"data"`
}

type BusinessProfileUpdateRequestDoc struct {
	BizName  string `json:"bizName" example:"Glow Beauty Studio"`
	BizType  string `json:"bizType" example:"salon_kecantikan"`
	BizSlug  string `json:"bizSlug" example:"glow-beauty-studio"`
	Timezone string `json:"timezone" example:"Asia/Jakarta"`
	Country  string `json:"country" example:"ID"`
}

type BusinessWhatsAppUpdateRequestDoc struct {
	OwnerWA         string `json:"ownerWa" example:"6281234567890"`
	WANum           string `json:"waNum" example:"6289876543210"`
	MetaWABAID      string `json:"metaWabaId" example:"1302223415207824"`
	MetaAccessToken string `json:"metaAccessToken" example:"EAALhDPG71UIBA..."`
}

type BusinessSettingsUpdateRequestDoc struct {
	AutomationEnabled      *bool  `json:"automationEnabled" example:"true"`
	DefaultInterval        *int   `json:"defaultInterval" example:"30"`
	SendTime               string `json:"sendTime" example:"09:00"`
	Timezone               string `json:"timezone" example:"Asia/Jakarta"`
	BillingNotifLow        *bool  `json:"billingNotifLow" example:"true"`
	BillingNotifCritical   *bool  `json:"billingNotifCritical" example:"true"`
	BillingNotifSubLow     *bool  `json:"billingNotifSubLow" example:"true"`
	BillingNotifPreRenewal *bool  `json:"billingNotifPreRenewal" example:"true"`
	AutoTopup              *bool  `json:"autoTopup" example:"false"`
	AutoTopupThreshold     *int   `json:"autoTopupThreshold" example:"10"`
	AutoTopupPackage       string `json:"autoTopupPackage" example:"p1"`
}

type CategoryItemDoc struct {
	ID             string `json:"id" example:"defcat-facial-treatment"`
	Name           string `json:"name" example:"Facial Treatment"`
	Icon           string `json:"icon" example:"💆"`
	Interval       int    `json:"interval" example:"30"`
	TemplateID     string `json:"templateId" example:"tpl-a"`
	TemplateBody   string `json:"templateBody" example:"Hai [nama], sudah sebulan sejak Facial Treatment terakhir kamu di [bisnis]."`
	MetaTemplateID string `json:"metaTemplateId,omitempty" example:"1744775703359541"`
	IsEnabled      bool   `json:"isEnabled" example:"true"`
}

type CategoryListResponseDoc struct {
	Error bool              `json:"error" example:"false"`
	Data  []CategoryItemDoc `json:"data"`
}

type DefaultCategoryItemDoc struct {
	ID           string `json:"id" example:"defcat-facial-treatment"`
	Name         string `json:"name" example:"Facial Treatment"`
	Icon         string `json:"icon" example:"💆"`
	Interval     int    `json:"interval" example:"30"`
	TemplateID   string `json:"templateId" example:"tpl-a"`
	TemplateBody string `json:"templateBody" example:"Halo 1! Sudah 2 hari sejak 3 terakhir kamu di 4. Yuk balik lagi — kami tunggu! 😊"`
	ExampleBody  string `json:"exampleBody" example:"[\"Pelanggan\",\"interval\",\"service\",\"business\"]"`
}

type DefaultCategoryListResponseDoc struct {
	Error bool                     `json:"error" example:"false"`
	Data  []DefaultCategoryItemDoc `json:"data"`
}

type CategoryCreateRequestDoc struct {
	Name              string `json:"name" example:"Hair Treatment"`
	Icon              string `json:"icon" example:"💇"`
	Interval          int    `json:"interval" example:"45"`
	DefaultCategoryID string `json:"defaultCategoryId" example:"defcat-hair-treatment"`
}

type CategoryCreateDataDoc struct {
	OK             bool   `json:"ok" example:"true"`
	CategoryID     string `json:"categoryId" example:"0ba9d96c-62e6-4ac8-bc77-b0228066f3ff"`
	MetaTemplateID string `json:"metaTemplateId" example:"1744775703359541"`
	MetaStatus     string `json:"metaStatus" example:"PENDING"`
}

type CategoryCreateResponseDoc struct {
	Error bool                  `json:"error" example:"false"`
	Data  CategoryCreateDataDoc `json:"data"`
}

type CategoryUpdateRequestDoc struct {
	Name         *string `json:"name" example:"Hair Treatment Premium"`
	Icon         *string `json:"icon" example:"💇"`
	Interval     *int    `json:"interval" example:"60"`
	TemplateID   *string `json:"templateId" example:"tpl-d"`
	TemplateBody *string `json:"templateBody" example:"Hai [nama], waktunya hair treatment premium lagi di [bisnis]."`
	IsEnabled    *bool   `json:"isEnabled" example:"true"`
}

type CustomerServiceInputDoc struct {
	CategoryID string `json:"categoryId" example:"defcat-facial-treatment"`
	Date       string `json:"date" example:"2026-04-12T09:00:00Z"`
}

type CustomerCreateRequestDoc struct {
	Name        string                    `json:"name" example:"Anisa Putri"`
	PhoneNumber string                    `json:"phoneNumber" example:"6281234567890"`
	Via         string                    `json:"via" example:"manual"`
	Services    []CustomerServiceInputDoc `json:"services"`
}

type CustomerUpdateRequestDoc struct {
	Name        *string                   `json:"name" example:"Anisa Putri"`
	PhoneNumber *string                   `json:"phoneNumber" example:"6281234567890"`
	Via         *string                   `json:"via" example:"manual"`
	Services    []CustomerServiceInputDoc `json:"services"`
}

type CustomerServiceDoc struct {
	Name   string `json:"name" example:"Facial Treatment"`
	Icon   string `json:"icon" example:"💆"`
	Date   string `json:"date" example:"2026-04-12T09:00:00Z"`
	Days   int    `json:"days" example:"30"`
	Status string `json:"status" example:"aktif"`
}

type CustomerListItemDoc struct {
	ID          string               `json:"id" example:"c9e53f5d-6dc7-4db5-8504-3625e0d737f5"`
	Name        string               `json:"name" example:"Anisa Putri"`
	PhoneNumber string               `json:"phoneNumber" example:"6281234567890"`
	Via         string               `json:"via" example:"manual"`
	Status      string               `json:"status" example:"aktif"`
	OverdueDays int                  `json:"overdueDays" example:"0"`
	CreatedAt   string               `json:"createdAt" example:"2026-04-12T08:30:00Z"`
	Services    []CustomerServiceDoc `json:"services"`
}

type CustomerPaginationDoc struct {
	Page       int  `json:"page" example:"1"`
	Limit      int  `json:"limit" example:"20"`
	Total      int  `json:"total" example:"57"`
	TotalPages int  `json:"totalPages" example:"3"`
	HasNext    bool `json:"hasNext" example:"true"`
	HasPrev    bool `json:"hasPrev" example:"false"`
}

type CustomerStatusCountDoc struct {
	Semua     int `json:"semua" example:"57"`
	Aktif     int `json:"aktif" example:"25"`
	Mendekati int `json:"mendekati" example:"20"`
	Hilang    int `json:"hilang" example:"12"`
}

type CustomerListResponseDoc struct {
	Error       bool                   `json:"error" example:"false"`
	Data        []CustomerListItemDoc  `json:"data"`
	Pagination  CustomerPaginationDoc  `json:"pagination"`
	StatusCount CustomerStatusCountDoc `json:"statusCount"`
}

type VisitRequestDoc struct {
	CustomerID          string   `json:"customerId" example:"c9e53f5d-6dc7-4db5-8504-3625e0d737f5"`
	CustomerName        string   `json:"customerName" example:"Anisa Putri"`
	CustomerPhoneNumber string   `json:"customerPhoneNumber,omitempty" example:"6281234567890"`
	Date                string   `json:"date" example:"2026-04-12T09:00:00Z"`
	CategoryIDs         []string `json:"categoryIds" example:"[\"defcat-facial-treatment\"]"`
}

type CheckinLookupRequestDoc struct {
	PhoneNumber string `json:"phoneNumber" example:"6281234567890"`
}

type CheckinLookupCustomerDoc struct {
	ID          string               `json:"id" example:"c9e53f5d-6dc7-4db5-8504-3625e0d737f5"`
	Name        string               `json:"name" example:"Anisa Putri"`
	PhoneNumber string               `json:"phoneNumber" example:"6281234567890"`
	Via         string               `json:"via" example:"manual"`
	Services    []CustomerServiceDoc `json:"services"`
}

type CheckinLookupDataDoc struct {
	Found    bool                      `json:"found" example:"true"`
	Customer *CheckinLookupCustomerDoc `json:"customer,omitempty"`
}

type CheckinLookupResponseDoc struct {
	Error bool                 `json:"error" example:"false"`
	Data  CheckinLookupDataDoc `json:"data"`
}

type CheckinSubmitRequestDoc struct {
	PhoneNumber string   `json:"phoneNumber" example:"6281234567890"`
	Name        string   `json:"name" example:"Anisa Putri"`
	Date        string   `json:"date" example:"2026-04-12T09:00:00Z"`
	CategoryIDs []string `json:"categoryIds" example:"[\"defcat-facial-treatment\"]"`
}

type ReminderLogItemDoc struct {
	ID            string `json:"id" example:"rem-1"`
	BusinessID    string `json:"businessId" example:"biz-1"`
	CustomerID    string `json:"customerId" example:"cx-1"`
	CategoryID    string `json:"categoryId" example:"cat-1"`
	CxName        string `json:"cxName" example:"Anisa Putri"`
	SvcName       string `json:"svcName" example:"Facial Treatment"`
	ScheduledAt   string `json:"scheduledAt" example:"2026-04-12T09:00:00Z"`
	SentAt        string `json:"sentAt" example:"2026-04-12T09:00:05Z"`
	Status        string `json:"status" example:"terkirim"`
	Kredit        int    `json:"kredit" example:"1"`
	ErrorReason   string `json:"errorReason,omitempty" example:""`
	RetryCount    int    `json:"retryCount" example:"0"`
	MetaMessageID string `json:"metaMessageId,omitempty" example:"wamid.HBgM..."`
}

type ReminderLogResponseDoc struct {
	Error bool                 `json:"error" example:"false"`
	Data  []ReminderLogItemDoc `json:"data"`
}

type DashboardCreditSummaryDoc struct {
	Welcome      int `json:"welcome" example:"100"`
	Subscription int `json:"subscription" example:"0"`
	Topup        int `json:"topup" example:"0"`
	Total        int `json:"total" example:"100"`
}

type DashboardSummaryDataDoc struct {
	Date               string                    `json:"date" example:"2026-04-12"`
	Timezone           string                    `json:"timezone" example:"Asia/Jakarta"`
	TotalCustomers     int                       `json:"totalCustomers" example:"12"`
	PendingReminders   int                       `json:"pendingReminders" example:"4"`
	SentToday          int                       `json:"sentToday" example:"8"`
	FailedToday        int                       `json:"failedToday" example:"1"`
	Credits            DashboardCreditSummaryDoc `json:"credits"`
	TrialEndsAt        string                    `json:"trialEndsAt" example:"2026-05-02T00:00:00Z"`
	SubscriptionStatus string                    `json:"subscriptionStatus" example:"none"`
}

type DashboardSummaryResponseDoc struct {
	Error bool                    `json:"error" example:"false"`
	Data  DashboardSummaryDataDoc `json:"data"`
}

type BillingTierDoc struct {
	Price   int `json:"price" example:"250000"`
	Credits int `json:"credits" example:"300"`
}

type BillingPlanConfigDoc struct {
	FreeBonus  int              `json:"freeBonus" example:"100"`
	SubCredits int              `json:"subCredits" example:"250"`
	SubPrice   int              `json:"subPrice" example:"250000"`
	TopupPrice int              `json:"topupPrice" example:"1000"`
	Tiers      []BillingTierDoc `json:"tiers"`
}

type BillingSummaryDataDoc struct {
	Plan               string               `json:"plan" example:"free"`
	SubscriptionStatus string               `json:"subscriptionStatus" example:"none"`
	TrialEndsAt        string               `json:"trialEndsAt" example:"2026-05-02T00:00:00Z"`
	SubscriptionEndsAt string               `json:"subscriptionEndsAt,omitempty" example:""`
	WelcomeCreditsLeft int                  `json:"welcomeCreditsLeft" example:"100"`
	SubCreditsLeft     int                  `json:"subCreditsLeft" example:"0"`
	TopupCreditsLeft   int                  `json:"topupCreditsLeft" example:"0"`
	SubCreditsMax      int                  `json:"subCreditsMax" example:"250"`
	RemLeft            int                  `json:"remLeft" example:"100"`
	PlanConfig         BillingPlanConfigDoc `json:"planConfig"`
}

type BillingSummaryResponseDoc struct {
	Error bool                  `json:"error" example:"false"`
	Data  BillingSummaryDataDoc `json:"data"`
}

type BillingHistoryItemDoc struct {
	ID           string `json:"id" example:"tx-1"`
	Type         string `json:"type" example:"welcome"`
	Label        string `json:"label" example:"Welcome Bonus"`
	Delta        int    `json:"delta" example:"100"`
	BalanceAfter int    `json:"balanceAfter" example:"100"`
	Note         string `json:"note" example:"Bonus kredit pendaftaran"`
	CreatedAt    string `json:"createdAt" example:"2026-04-12T09:00:00Z"`
}

type BillingHistoryResponseDoc struct {
	Error bool                    `json:"error" example:"false"`
	Data  []BillingHistoryItemDoc `json:"data"`
}

type BillingCheckoutRequestDoc struct {
	PackageID string `json:"packageId" example:"p1"`
}

type BillingCheckoutDataDoc struct {
	OrderID     string `json:"orderId" example:"3cdfa11e-ae03-4ddd-93c2-bcb6219f9859"`
	ExternalID  string `json:"externalId" example:"topup-6f2f4ea7-45e7-4f8f-af3a-bfcba1f3f8db"`
	InvoiceID   string `json:"invoiceId" example:"64f2c6ab-0194-4d33-81e2-42847eb9190d"`
	CheckoutURL string `json:"checkoutUrl" example:"https://checkout.xendit.co/web/64f2c6ab..."`
	Status      string `json:"status" example:"PENDING"`
}

type BillingCheckoutResponseDoc struct {
	Error bool                   `json:"error" example:"false"`
	Data  BillingCheckoutDataDoc `json:"data"`
}

type XenditWebhookPayloadDoc struct {
	ID         string `json:"id" example:"64f2c6ab-0194-4d33-81e2-42847eb9190d"`
	ExternalID string `json:"external_id" example:"topup-6f2f4ea7-45e7-4f8f-af3a-bfcba1f3f8db"`
	Status     string `json:"status" example:"PAID"`
	PaidAt     string `json:"paid_at" example:"2026-04-12T10:20:00Z"`
}

type MetaWebhookPayloadDoc struct {
	Object string `json:"object" example:"whatsapp_business_account"`
	Entry  []struct {
		ID      string `json:"id" example:"1302223415207824"`
		Changes []struct {
			Field string `json:"field" example:"messages"`
		} `json:"changes"`
	} `json:"entry"`
}

type InternalPlanConfigResponseDoc struct {
	Error bool `json:"error" example:"false"`
	Data  struct {
		BusinessID string `json:"businessId" example:"biz-1"`
		FreeBonus  int    `json:"freeBonus" example:"100"`
		SubCredits int    `json:"subCredits" example:"250"`
		SubPrice   int    `json:"subPrice" example:"250000"`
		TopupPrice int    `json:"topupPrice" example:"1000"`
		Tier1Price int    `json:"tier1Price" example:"250000"`
		Tier1Creds int    `json:"tier1Credits" example:"300"`
		Tier2Price int    `json:"tier2Price" example:"500000"`
		Tier2Creds int    `json:"tier2Credits" example:"625"`
		Tier3Price int    `json:"tier3Price" example:"1000000"`
		Tier3Creds int    `json:"tier3Credits" example:"1500"`
	} `json:"data"`
}

type InternalPlanConfigUpdateRequestDoc struct {
	FreeBonus    int `json:"freeBonus" example:"100"`
	SubCredits   int `json:"subCredits" example:"250"`
	SubPrice     int `json:"subPrice" example:"250000"`
	TopupPrice   int `json:"topupPrice" example:"1000"`
	Tier1Price   int `json:"tier1Price" example:"250000"`
	Tier1Credits int `json:"tier1Credits" example:"300"`
	Tier2Price   int `json:"tier2Price" example:"500000"`
	Tier2Credits int `json:"tier2Credits" example:"625"`
	Tier3Price   int `json:"tier3Price" example:"1000000"`
	Tier3Credits int `json:"tier3Credits" example:"1500"`
}

type InternalAdminLoginRequestDoc struct {
	Email    string `json:"email" example:"admin@getstarvio.com"`
	Password string `json:"password" example:"secret-password"`
}

type InternalAdminLoginResponseDoc struct {
	Error bool `json:"error" example:"false"`
	Data  struct {
		AccessToken  string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIs..."`
		RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIs..."`
		Admin        struct {
			ID    string `json:"id" example:"d4f89973-19a8-4b53-a736-96f9c4fc36bf"`
			Name  string `json:"name" example:"System Admin"`
			Email string `json:"email" example:"admin@getstarvio.com"`
		} `json:"admin"`
	} `json:"data"`
}

type InternalAdminRefreshRequestDoc struct {
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIs..."`
}

type InternalAdminRefreshResponseDoc struct {
	Error bool `json:"error" example:"false"`
	Data  struct {
		AccessToken  string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIs..."`
		RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIs..."`
	} `json:"data"`
}

type InternalAdminLogoutRequestDoc struct {
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// healthzDoc godoc
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /healthz [get]
func healthzDoc() {}

// authGoogleLoginDoc godoc
// @Summary Login with Google
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body AuthGoogleLoginRequestDoc true "Google login payload"
// @Success 200 {object} AuthTokenResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/auth/google/login [post]
func authGoogleLoginDoc() {}

// authRefreshDoc godoc
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body AuthRefreshRequestDoc true "Refresh token payload"
// @Success 200 {object} AuthTokenResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/auth/refresh [post]
func authRefreshDoc() {}

// authLogoutDoc godoc
// @Summary Logout
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body AuthLogoutRequestDoc true "Logout payload"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Router /v1/auth/logout [post]
func authLogoutDoc() {}

// meBootstrapDoc godoc
// @Summary Get bootstrap payload
// @Tags business
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/me/bootstrap [get]
func meBootstrapDoc() {}

// businessProfileUpdateDoc godoc
// @Summary Update business profile
// @Tags business
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body BusinessProfileUpdateRequestDoc true "Business profile payload"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/business/profile [put]
func businessProfileUpdateDoc() {}

// businessWhatsAppUpdateDoc godoc
// @Summary Update business WhatsApp fields
// @Tags business
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body BusinessWhatsAppUpdateRequestDoc true "Business WhatsApp payload"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/business/whatsapp [put]
func businessWhatsAppUpdateDoc() {}

// businessSettingsUpdateDoc godoc
// @Summary Update business settings
// @Tags business
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body BusinessSettingsUpdateRequestDoc true "Business settings payload"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/business/settings [put]
func businessSettingsUpdateDoc() {}

// categoriesListDoc godoc
// @Summary List categories
// @Tags category
// @Security BearerAuth
// @Produce json
// @Success 200 {object} CategoryListResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/categories [get]
func categoriesListDoc() {}

// defaultCategoriesListDoc godoc
// @Summary List default categories
// @Tags category
// @Security BearerAuth
// @Produce json
// @Success 200 {object} DefaultCategoryListResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/default-categories [get]
func defaultCategoriesListDoc() {}

// categoriesCreateDoc godoc
// @Summary Create category
// @Description Creates WhatsApp template in Meta first. Business must have `metaWabaId` and `metaAccessToken` set via `/v1/business/whatsapp`. Returns error when Meta status is REJECTED.
// @Tags category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body CategoryCreateRequestDoc true "Category payload"
// @Success 201 {object} CategoryCreateResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/categories [post]
func categoriesCreateDoc() {}

// categoriesUpdateDoc godoc
// @Summary Update category
// @Tags category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param payload body CategoryUpdateRequestDoc true "Category payload"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/categories/{id} [patch]
func categoriesUpdateDoc() {}

// categoriesDeleteDoc godoc
// @Summary Delete category
// @Tags category
// @Security BearerAuth
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/categories/{id} [delete]
func categoriesDeleteDoc() {}

// customersListDoc godoc
// @Summary List customers
// @Tags customer
// @Security BearerAuth
// @Produce json
// @Param q query string false "Search query"
// @Param status query string false "Customer status"
// @Param sort query string false "Sort mode" default(urgent)
// @Param date query string false "Reference date (RFC3339 or YYYY-MM-DD)" example(2026-04-14)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page (max 100)" default(20)
// @Success 200 {object} CustomerListResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/customers [get]
func customersListDoc() {}

// customersCreateDoc godoc
// @Summary Create customer
// @Tags customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body CustomerCreateRequestDoc true "Customer payload"
// @Success 201 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/customers [post]
func customersCreateDoc() {}

// customersUpdateDoc godoc
// @Summary Update customer
// @Tags customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param payload body CustomerUpdateRequestDoc true "Customer payload"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/customers/{id} [patch]
func customersUpdateDoc() {}

// customersDeleteDoc godoc
// @Summary Delete customer
// @Tags customer
// @Security BearerAuth
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/customers/{id} [delete]
func customersDeleteDoc() {}

// visitsCreateDoc godoc
// @Summary Create manual visit
// @Description Supports backdate up to 7 days.
// @Tags customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body VisitRequestDoc true "Visit payload"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/visits [post]
func visitsCreateDoc() {}

// checkinLookupDoc godoc
// @Summary Lookup customer for check-in
// @Tags customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body CheckinLookupRequestDoc true "Lookup payload"
// @Success 200 {object} CheckinLookupResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/checkin/lookup [post]
func checkinLookupDoc() {}

// checkinSubmitDoc godoc
// @Summary Submit check-in
// @Tags customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body CheckinSubmitRequestDoc true "Check-in payload"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/checkin/submit [post]
func checkinSubmitDoc() {}

// remindersLogDoc godoc
// @Summary List reminder logs
// @Tags reminder
// @Security BearerAuth
// @Produce json
// @Param status query string false "Reminder status"
// @Param limit query int false "Limit" default(200)
// @Success 200 {object} ReminderLogResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/reminders/log [get]
func remindersLogDoc() {}

// remindersRetryDoc godoc
// @Summary Retry reminder
// @Tags reminder
// @Security BearerAuth
// @Produce json
// @Param id path string true "Reminder ID"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/reminders/{id}/retry [post]
func remindersRetryDoc() {}

// dashboardSummaryDoc godoc
// @Summary Dashboard summary
// @Tags reminder
// @Security BearerAuth
// @Produce json
// @Success 200 {object} DashboardSummaryResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/dashboard/summary [get]
func dashboardSummaryDoc() {}

// billingSummaryDoc godoc
// @Summary Billing summary
// @Tags billing
// @Security BearerAuth
// @Produce json
// @Success 200 {object} BillingSummaryResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/billing/summary [get]
func billingSummaryDoc() {}

// billingHistoryDoc godoc
// @Summary Billing history
// @Tags billing
// @Security BearerAuth
// @Produce json
// @Success 200 {object} BillingHistoryResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/billing/history [get]
func billingHistoryDoc() {}

// billingSubscriptionActivateDoc godoc
// @Summary Activate subscription
// @Description Simulated activation in v1.
// @Tags billing
// @Security BearerAuth
// @Produce json
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/billing/subscription/activate [post]
func billingSubscriptionActivateDoc() {}

// billingSubscriptionCancelDoc godoc
// @Summary Cancel subscription
// @Tags billing
// @Security BearerAuth
// @Produce json
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/billing/subscription/cancel [post]
func billingSubscriptionCancelDoc() {}

// billingTopupCheckoutDoc godoc
// @Summary Create top-up checkout (Xendit)
// @Tags billing
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body BillingCheckoutRequestDoc true "Checkout payload"
// @Success 200 {object} BillingCheckoutResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/billing/topup/checkout [post]
func billingTopupCheckoutDoc() {}

// webhooksXenditDoc godoc
// @Summary Xendit webhook
// @Tags webhook
// @Accept json
// @Produce json
// @Param X-Callback-Token header string false "Xendit callback token"
// @Param payload body XenditWebhookPayloadDoc true "Xendit webhook payload"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/webhooks/xendit [post]
func webhooksXenditDoc() {}

// webhooksMetaDoc godoc
// @Summary Meta webhook event receiver
// @Tags webhook
// @Accept json
// @Produce json
// @Param X-Hub-Signature-256 header string false "Meta payload signature"
// @Param payload body MetaWebhookPayloadDoc true "Meta webhook payload"
// @Success 200 {string} string "EVENT_RECEIVED"
// @Failure 401 {string} string "unauthorized"
// @Failure 400 {string} string "bad request"
// @Router /v1/webhooks/meta [post]
func webhooksMetaDoc() {}

// webhooksMetaVerifyDoc godoc
// @Summary Meta webhook verification
// @Tags webhook
// @Produce plain
// @Param hub.mode query string true "subscribe"
// @Param hub.verify_token query string true "Verification token"
// @Param hub.challenge query string true "Challenge value"
// @Success 200 {string} string "hub.challenge"
// @Failure 400 {string} string "bad request"
// @Failure 401 {string} string "unauthorized"
// @Router /v1/webhooks/meta [get]
func webhooksMetaVerifyDoc() {}

// internalAdminLoginDoc godoc
// @Summary Internal admin login
// @Tags internal
// @Accept json
// @Produce json
// @Param payload body InternalAdminLoginRequestDoc true "Internal admin login payload"
// @Success 200 {object} InternalAdminLoginResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/internal/auth/login [post]
func internalAdminLoginDoc() {}

// internalAdminRefreshDoc godoc
// @Summary Internal admin refresh token
// @Tags internal
// @Accept json
// @Produce json
// @Param payload body InternalAdminRefreshRequestDoc true "Internal admin refresh payload"
// @Success 200 {object} InternalAdminRefreshResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/internal/auth/refresh [post]
func internalAdminRefreshDoc() {}

// internalAdminLogoutDoc godoc
// @Summary Internal admin logout
// @Tags internal
// @Accept json
// @Produce json
// @Param payload body InternalAdminLogoutRequestDoc true "Internal admin logout payload"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/internal/auth/logout [post]
func internalAdminLogoutDoc() {}

// internalPlanConfigGetDoc godoc
// @Summary Get internal plan config
// @Tags internal
// @Security BearerAuth
// @Produce json
// @Success 200 {object} InternalPlanConfigResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/internal/plan-config [get]
func internalPlanConfigGetDoc() {}

// internalPlanConfigPutDoc godoc
// @Summary Update internal plan config
// @Tags internal
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body InternalPlanConfigUpdateRequestDoc true "Plan config payload"
// @Success 200 {object} OKResponseDoc
// @Failure 400 {object} ErrorResponseDoc
// @Failure 401 {object} ErrorResponseDoc
// @Router /v1/internal/plan-config [put]
func internalPlanConfigPutDoc() {}
