package routes

// SwaggerEnvelope is the common API response shape.
type SwaggerEnvelope struct {
	Error   bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// HealthResponse is the health check payload.
type HealthResponse struct {
	OK      bool   `json:"ok"`
	Service string `json:"service"`
	Time    string `json:"time"`
}

// GenericBody is a flexible request payload placeholder for docs.
type GenericBody map[string]interface{}

// AuthTokenBody is login/refresh response payload.
type AuthTokenBody struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userId"`
}

// InternalTokenHeader documents internal API token header.
type InternalTokenHeader struct {
	XInternalToken string `json:"X-Internal-Token"`
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
// @Param payload body GenericBody true "Google login payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/auth/google/login [post]
func authGoogleLoginDoc() {}

// authRefreshDoc godoc
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body GenericBody true "Refresh token payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/auth/refresh [post]
func authRefreshDoc() {}

// authLogoutDoc godoc
// @Summary Logout
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body GenericBody true "Logout payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Router /v1/auth/logout [post]
func authLogoutDoc() {}

// meBootstrapDoc godoc
// @Summary Get bootstrap payload
// @Tags business
// @Security BearerAuth
// @Produce json
// @Success 200 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/me/bootstrap [get]
func meBootstrapDoc() {}

// businessProfileUpdateDoc godoc
// @Summary Update business profile
// @Tags business
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body GenericBody true "Business profile payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/business/profile [put]
func businessProfileUpdateDoc() {}

// businessWhatsAppUpdateDoc godoc
// @Summary Update business WhatsApp fields
// @Tags business
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body GenericBody true "Business WhatsApp payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/business/whatsapp [put]
func businessWhatsAppUpdateDoc() {}

// businessSettingsUpdateDoc godoc
// @Summary Update business settings
// @Tags business
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body GenericBody true "Business settings payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/business/settings [put]
func businessSettingsUpdateDoc() {}

// categoriesListDoc godoc
// @Summary List categories
// @Tags category
// @Security BearerAuth
// @Produce json
// @Success 200 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/categories [get]
func categoriesListDoc() {}

// categoriesCreateDoc godoc
// @Summary Create category
// @Tags category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body GenericBody true "Category payload"
// @Success 201 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/categories [post]
func categoriesCreateDoc() {}

// categoriesUpdateDoc godoc
// @Summary Update category
// @Tags category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param payload body GenericBody true "Category payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/categories/{id} [patch]
func categoriesUpdateDoc() {}

// categoriesDeleteDoc godoc
// @Summary Delete category
// @Tags category
// @Security BearerAuth
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
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
// @Success 200 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/customers [get]
func customersListDoc() {}

// customersCreateDoc godoc
// @Summary Create customer
// @Tags customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body GenericBody true "Customer payload"
// @Success 201 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/customers [post]
func customersCreateDoc() {}

// customersUpdateDoc godoc
// @Summary Update customer
// @Tags customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param payload body GenericBody true "Customer payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/customers/{id} [patch]
func customersUpdateDoc() {}

// customersDeleteDoc godoc
// @Summary Delete customer
// @Tags customer
// @Security BearerAuth
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/customers/{id} [delete]
func customersDeleteDoc() {}

// visitsCreateDoc godoc
// @Summary Create manual visit
// @Description Supports backdate up to 7 days.
// @Tags customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body GenericBody true "Visit payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/visits [post]
func visitsCreateDoc() {}

// checkinLookupDoc godoc
// @Summary Lookup customer for check-in
// @Tags customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body GenericBody true "Lookup payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/checkin/lookup [post]
func checkinLookupDoc() {}

// checkinSubmitDoc godoc
// @Summary Submit check-in
// @Tags customer
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body GenericBody true "Check-in payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/checkin/submit [post]
func checkinSubmitDoc() {}

// remindersLogDoc godoc
// @Summary List reminder logs
// @Tags reminder
// @Security BearerAuth
// @Produce json
// @Param status query string false "Reminder status"
// @Param limit query int false "Limit" default(200)
// @Success 200 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/reminders/log [get]
func remindersLogDoc() {}

// remindersRetryDoc godoc
// @Summary Retry reminder
// @Tags reminder
// @Security BearerAuth
// @Produce json
// @Param id path string true "Reminder ID"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/reminders/{id}/retry [post]
func remindersRetryDoc() {}

// dashboardSummaryDoc godoc
// @Summary Dashboard summary
// @Tags reminder
// @Security BearerAuth
// @Produce json
// @Success 200 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/dashboard/summary [get]
func dashboardSummaryDoc() {}

// billingSummaryDoc godoc
// @Summary Billing summary
// @Tags billing
// @Security BearerAuth
// @Produce json
// @Success 200 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/billing/summary [get]
func billingSummaryDoc() {}

// billingHistoryDoc godoc
// @Summary Billing history
// @Tags billing
// @Security BearerAuth
// @Produce json
// @Success 200 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/billing/history [get]
func billingHistoryDoc() {}

// billingSubscriptionActivateDoc godoc
// @Summary Activate subscription
// @Description Simulated activation in v1.
// @Tags billing
// @Security BearerAuth
// @Produce json
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/billing/subscription/activate [post]
func billingSubscriptionActivateDoc() {}

// billingSubscriptionCancelDoc godoc
// @Summary Cancel subscription
// @Tags billing
// @Security BearerAuth
// @Produce json
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/billing/subscription/cancel [post]
func billingSubscriptionCancelDoc() {}

// billingTopupCheckoutDoc godoc
// @Summary Create top-up checkout (Xendit)
// @Tags billing
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body GenericBody true "Checkout payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/billing/topup/checkout [post]
func billingTopupCheckoutDoc() {}

// webhooksXenditDoc godoc
// @Summary Xendit webhook
// @Tags webhook
// @Accept json
// @Produce json
// @Param X-Callback-Token header string false "Xendit callback token"
// @Param payload body GenericBody true "Xendit webhook payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/webhooks/xendit [post]
func webhooksXenditDoc() {}

// webhooksMetaDoc godoc
// @Summary Meta webhook
// @Tags webhook
// @Accept json
// @Produce json
// @Param payload body GenericBody true "Meta webhook payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Router /v1/webhooks/meta [post]
func webhooksMetaDoc() {}

// internalPlanConfigGetDoc godoc
// @Summary Get internal plan config
// @Tags internal
// @Produce json
// @Param X-Internal-Token header string true "Internal API token"
// @Success 200 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/internal/plan-config [get]
func internalPlanConfigGetDoc() {}

// internalPlanConfigPutDoc godoc
// @Summary Update internal plan config
// @Tags internal
// @Accept json
// @Produce json
// @Param X-Internal-Token header string true "Internal API token"
// @Param payload body GenericBody true "Plan config payload"
// @Success 200 {object} SwaggerEnvelope
// @Failure 400 {object} SwaggerEnvelope
// @Failure 401 {object} SwaggerEnvelope
// @Router /v1/internal/plan-config [put]
func internalPlanConfigPutDoc() {}
