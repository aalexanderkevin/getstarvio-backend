package models

import "time"

const (
	ReminderStatusPending = "pending"
	ReminderStatusSent    = "terkirim"
	ReminderStatusFailed  = "gagal"

	SubscriptionStatusNone      = "none"
	SubscriptionStatusActive    = "active"
	SubscriptionStatusCancelled = "cancelled"
)

type User struct {
	ID        string    `gorm:"column:id;primaryKey"`
	GoogleSub string    `gorm:"column:google_sub;uniqueIndex;not null"`
	Email     string    `gorm:"column:email;uniqueIndex;not null"`
	Name      string    `gorm:"column:name;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string { return "users" }

type Business struct {
	ID              string    `gorm:"column:id;primaryKey"`
	UserID          string    `gorm:"column:user_id;uniqueIndex;not null"`
	BizName         string    `gorm:"column:biz_name;not null"`
	BizType         string    `gorm:"column:biz_type;not null"`
	BizSlug         string    `gorm:"column:biz_slug;uniqueIndex;not null"`
	AdminName       string    `gorm:"column:admin_name;not null"`
	AdminEmail      string    `gorm:"column:admin_email;not null"`
	OwnerWA         string    `gorm:"column:owner_wa"`
	WANum           string    `gorm:"column:wa_num"`
	MetaWABAID      string    `gorm:"column:meta_waba_id"`
	MetaAccessToken string    `gorm:"column:meta_access_token"`
	Timezone        string    `gorm:"column:timezone;not null"`
	Country         string    `gorm:"column:country;not null"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Business) TableName() string { return "businesses" }

type BusinessSettings struct {
	ID                   string    `gorm:"column:id;primaryKey"`
	BusinessID           string    `gorm:"column:business_id;uniqueIndex;not null"`
	AutomationEnabled    bool      `gorm:"column:automation_enabled;default:true"`
	DefaultInterval      int       `gorm:"column:default_interval;default:30"`
	SendTime             string    `gorm:"column:send_time;default:'09:00'"`
	Timezone             string    `gorm:"column:timezone;not null"`
	BillingNotifLow      bool      `gorm:"column:billing_notif_low;default:true"`
	BillingNotifCritical bool      `gorm:"column:billing_notif_critical;default:true"`
	BillingNotifSubLow   bool      `gorm:"column:billing_notif_sub_low;default:true"`
	BillingNotifPreRenew bool      `gorm:"column:billing_notif_pre_renewal;default:true"`
	AutoTopupEnabled     bool      `gorm:"column:auto_topup_enabled;default:false"`
	AutoTopupThreshold   int       `gorm:"column:auto_topup_threshold;default:10"`
	AutoTopupPackageID   string    `gorm:"column:auto_topup_package_id;default:'p1'"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (BusinessSettings) TableName() string { return "business_settings" }

type Category struct {
	ID             string    `gorm:"column:id;primaryKey"`
	BusinessID     string    `gorm:"column:business_id;index;not null"`
	Name           string    `gorm:"column:name;not null"`
	Icon           string    `gorm:"column:icon;not null"`
	IntervalDays   int       `gorm:"column:interval_days;not null"`
	TemplateID     string    `gorm:"column:template_id;not null"`
	TemplateBody   string    `gorm:"column:template_body;not null"`
	MetaTemplateID string    `gorm:"column:meta_template_id"`
	IsEnabled      bool      `gorm:"column:is_enabled;default:true"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Category) TableName() string { return "categories" }

type DefaultCategory struct {
	ID           string    `gorm:"column:id;primaryKey"`
	Name         string    `gorm:"column:name;not null"`
	Category     string    `gorm:"column:category;not null;default:'UTILITY'"`
	Status       string    `gorm:"column:status;not null;default:'PENDING'"`
	Icon         *string   `gorm:"column:icon"`
	IntervalDays *int      `gorm:"column:interval_days"`
	TemplateID   string    `gorm:"column:template_id;not null"`
	TemplateBody string    `gorm:"column:template_body;not null"`
	ExampleBody  string    `gorm:"column:example_body;not null"`
	IsActive     bool      `gorm:"column:is_active;default:true"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (DefaultCategory) TableName() string { return "default_categories" }

type Customer struct {
	ID          string    `gorm:"column:id;primaryKey"`
	BusinessID  string    `gorm:"column:business_id;index;not null"`
	Name        string    `gorm:"column:name;not null"`
	PhoneNumber string    `gorm:"column:phone_number;index;not null"`
	Via         string    `gorm:"column:via;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Customer) TableName() string { return "customers" }

type CustomerService struct {
	ID           string    `gorm:"column:id;primaryKey"`
	CustomerID   string    `gorm:"column:customer_id;index;not null"`
	CategoryID   string    `gorm:"column:category_id;index"`
	ServiceName  string    `gorm:"column:service_name;->"`
	ServiceIcon  string    `gorm:"column:service_icon;->"`
	LastVisitAt  time.Time `gorm:"column:last_visit_at;not null"`
	IntervalDays int       `gorm:"column:interval_days;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (CustomerService) TableName() string { return "customer_services" }

type Reminder struct {
	ID            string     `gorm:"column:id;primaryKey"`
	BusinessID    string     `gorm:"column:business_id;index;not null"`
	CustomerID    string     `gorm:"column:customer_id;index;not null"`
	CategoryID    string     `gorm:"column:category_id;index"`
	CxName        string     `gorm:"column:cx_name;not null"`
	SvcName       string     `gorm:"column:svc_name;not null"`
	ScheduledAt   time.Time  `gorm:"column:scheduled_at;index;not null"`
	SentAt        *time.Time `gorm:"column:sent_at"`
	Status        string     `gorm:"column:status;index;not null"`
	Kredit        int        `gorm:"column:kredit;not null;default:1"`
	ErrorReason   string     `gorm:"column:error_reason"`
	RetryCount    int        `gorm:"column:retry_count;default:0"`
	MetaMessageID string     `gorm:"column:meta_message_id"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (Reminder) TableName() string { return "reminders" }

type Wallet struct {
	ID                  string     `gorm:"column:id;primaryKey"`
	BusinessID          string     `gorm:"column:business_id;uniqueIndex;not null"`
	TrialStartedAt      time.Time  `gorm:"column:trial_started_at;not null"`
	TrialEndsAt         time.Time  `gorm:"column:trial_ends_at;not null"`
	SubscriptionStatus  string     `gorm:"column:subscription_status;not null;default:'none'"`
	SubscriptionStarted *time.Time `gorm:"column:subscription_started_at"`
	SubscriptionEnds    *time.Time `gorm:"column:subscription_ends_at"`
	WelcomeCreditsLeft  int        `gorm:"column:welcome_credits_left;not null;default:100"`
	SubCreditsLeft      int        `gorm:"column:sub_credits_left;not null;default:0"`
	TopupCreditsLeft    int        `gorm:"column:topup_credits_left;not null;default:0"`
	SubCreditsMax       int        `gorm:"column:sub_credits_max;not null;default:250"`
	CreatedAt           time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (Wallet) TableName() string { return "wallets" }

type BillingTransaction struct {
	ID           string    `gorm:"column:id;primaryKey"`
	BusinessID   string    `gorm:"column:business_id;index;not null"`
	Type         string    `gorm:"column:type;index;not null"`
	Label        string    `gorm:"column:label;not null"`
	Delta        int       `gorm:"column:delta;not null"`
	BalanceAfter int       `gorm:"column:balance_after;not null"`
	Note         string    `gorm:"column:note"`
	MetaJSON     string    `gorm:"column:meta_json"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (BillingTransaction) TableName() string { return "billing_transactions" }

type TopupOrder struct {
	ID          string     `gorm:"column:id;primaryKey"`
	BusinessID  string     `gorm:"column:business_id;index;not null"`
	ExternalID  string     `gorm:"column:external_id;uniqueIndex;not null"`
	InvoiceID   string     `gorm:"column:invoice_id;index"`
	PackageID   string     `gorm:"column:package_id;not null"`
	AmountIDR   int        `gorm:"column:amount_idr;not null"`
	Credits     int        `gorm:"column:credits;not null"`
	Status      string     `gorm:"column:status;index;not null"`
	CheckoutURL string     `gorm:"column:checkout_url"`
	PaidAt      *time.Time `gorm:"column:paid_at"`
	RawPayload  string     `gorm:"column:raw_payload"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (TopupOrder) TableName() string { return "topup_orders" }

type PlanConfig struct {
	ID           string    `gorm:"column:id;primaryKey"`
	BusinessID   string    `gorm:"column:business_id;uniqueIndex;not null"`
	FreeBonus    int       `gorm:"column:free_bonus;not null;default:100"`
	SubCredits   int       `gorm:"column:sub_credits;not null;default:250"`
	SubPrice     int       `gorm:"column:sub_price;not null;default:250000"`
	TopupPrice   int       `gorm:"column:topup_price;not null;default:1000"`
	Tier1Price   int       `gorm:"column:tier1_price;not null;default:250000"`
	Tier1Credits int       `gorm:"column:tier1_credits;not null;default:300"`
	Tier2Price   int       `gorm:"column:tier2_price;not null;default:500000"`
	Tier2Credits int       `gorm:"column:tier2_credits;not null;default:625"`
	Tier3Price   int       `gorm:"column:tier3_price;not null;default:1000000"`
	Tier3Credits int       `gorm:"column:tier3_credits;not null;default:1500"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (PlanConfig) TableName() string { return "plan_configs" }

type InternalAdmin struct {
	ID           string     `gorm:"column:id;primaryKey"`
	Name         string     `gorm:"column:name;not null"`
	Email        string     `gorm:"column:email;uniqueIndex;not null"`
	PasswordHash string     `gorm:"column:password_hash;not null"`
	IsActive     bool       `gorm:"column:is_active;not null;default:true"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (InternalAdmin) TableName() string { return "internal_admins" }

type InternalRefreshToken struct {
	ID        string     `gorm:"column:id;primaryKey"`
	AdminID   string     `gorm:"column:admin_id;index;not null"`
	TokenHash string     `gorm:"column:token_hash;uniqueIndex;not null"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (InternalRefreshToken) TableName() string { return "internal_refresh_tokens" }

type RefreshToken struct {
	ID        string     `gorm:"column:id;primaryKey"`
	UserID    string     `gorm:"column:user_id;index;not null"`
	TokenHash string     `gorm:"column:token_hash;uniqueIndex;not null"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (RefreshToken) TableName() string { return "refresh_tokens" }
