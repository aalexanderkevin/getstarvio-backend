CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  google_sub TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS businesses (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  biz_name TEXT NOT NULL,
  biz_type TEXT NOT NULL,
  biz_slug TEXT NOT NULL UNIQUE,
  admin_name TEXT NOT NULL,
  admin_email TEXT NOT NULL,
  owner_wa TEXT,
  wa_num TEXT,
  meta_waba_id TEXT,
  meta_access_token TEXT,
  timezone TEXT NOT NULL,
  country TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS business_settings (
  id TEXT PRIMARY KEY,
  business_id TEXT NOT NULL UNIQUE REFERENCES businesses(id) ON DELETE CASCADE,
  automation_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  default_interval INT NOT NULL DEFAULT 30,
  send_time TEXT NOT NULL DEFAULT '09:00',
  timezone TEXT NOT NULL,
  billing_notif_low BOOLEAN NOT NULL DEFAULT TRUE,
  billing_notif_critical BOOLEAN NOT NULL DEFAULT TRUE,
  billing_notif_sub_low BOOLEAN NOT NULL DEFAULT TRUE,
  billing_notif_pre_renewal BOOLEAN NOT NULL DEFAULT TRUE,
  auto_topup_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  auto_topup_threshold INT NOT NULL DEFAULT 10,
  auto_topup_package_id TEXT NOT NULL DEFAULT 'p1',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS default_categories (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  icon TEXT NOT NULL,
  interval_days INT NOT NULL,
  template_id TEXT NOT NULL,
  template_body TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_default_categories_is_active ON default_categories(is_active);

CREATE TABLE IF NOT EXISTS categories (
  id TEXT PRIMARY KEY,
  business_id TEXT NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  icon TEXT NOT NULL,
  interval_days INT NOT NULL,
  template_id TEXT NOT NULL,
  template_body TEXT NOT NULL,
  meta_template_id TEXT,
  is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_categories_business ON categories(business_id);

CREATE TABLE IF NOT EXISTS customers (
  id TEXT PRIMARY KEY,
  business_id TEXT NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  phone_number TEXT NOT NULL,
  via TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (business_id, phone_number)
);
CREATE INDEX IF NOT EXISTS idx_customers_business ON customers(business_id);
CREATE INDEX IF NOT EXISTS idx_customers_phone_number ON customers(phone_number);

CREATE TABLE IF NOT EXISTS customer_services (
  id TEXT PRIMARY KEY,
  customer_id TEXT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
  category_id TEXT REFERENCES categories(id) ON DELETE SET NULL,
  last_visit_at TIMESTAMPTZ NOT NULL,
  interval_days INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (customer_id, category_id)
);
CREATE INDEX IF NOT EXISTS idx_customer_services_customer ON customer_services(customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_services_category ON customer_services(category_id);

CREATE TABLE IF NOT EXISTS reminders (
  id TEXT PRIMARY KEY,
  business_id TEXT NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  customer_id TEXT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
  category_id TEXT REFERENCES categories(id) ON DELETE SET NULL,
  cx_name TEXT NOT NULL,
  svc_name TEXT NOT NULL,
  scheduled_at TIMESTAMPTZ NOT NULL,
  sent_at TIMESTAMPTZ,
  status TEXT NOT NULL,
  kredit INT NOT NULL DEFAULT 1,
  error_reason TEXT,
  retry_count INT NOT NULL DEFAULT 0,
  meta_message_id TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_reminders_business ON reminders(business_id);
CREATE INDEX IF NOT EXISTS idx_reminders_status ON reminders(status);
CREATE INDEX IF NOT EXISTS idx_reminders_due ON reminders(status, scheduled_at);
CREATE UNIQUE INDEX IF NOT EXISTS uq_reminders_identity ON reminders(business_id, customer_id, svc_name, scheduled_at);

CREATE TABLE IF NOT EXISTS wallets (
  id TEXT PRIMARY KEY,
  business_id TEXT NOT NULL UNIQUE REFERENCES businesses(id) ON DELETE CASCADE,
  trial_started_at TIMESTAMPTZ NOT NULL,
  trial_ends_at TIMESTAMPTZ NOT NULL,
  subscription_status TEXT NOT NULL DEFAULT 'none',
  subscription_started_at TIMESTAMPTZ,
  subscription_ends_at TIMESTAMPTZ,
  welcome_credits_left INT NOT NULL DEFAULT 100,
  sub_credits_left INT NOT NULL DEFAULT 0,
  topup_credits_left INT NOT NULL DEFAULT 0,
  sub_credits_max INT NOT NULL DEFAULT 250,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS billing_transactions (
  id TEXT PRIMARY KEY,
  business_id TEXT NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  type TEXT NOT NULL,
  label TEXT NOT NULL,
  delta INT NOT NULL,
  balance_after INT NOT NULL,
  note TEXT,
  meta_json TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_billing_transactions_business ON billing_transactions(business_id, created_at DESC);

CREATE TABLE IF NOT EXISTS topup_orders (
  id TEXT PRIMARY KEY,
  business_id TEXT NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  external_id TEXT NOT NULL UNIQUE,
  invoice_id TEXT,
  package_id TEXT NOT NULL,
  amount_idr INT NOT NULL,
  credits INT NOT NULL,
  status TEXT NOT NULL,
  checkout_url TEXT,
  paid_at TIMESTAMPTZ,
  raw_payload TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_topup_orders_business ON topup_orders(business_id);
CREATE INDEX IF NOT EXISTS idx_topup_orders_invoice ON topup_orders(invoice_id);

CREATE TABLE IF NOT EXISTS plan_configs (
  id TEXT PRIMARY KEY,
  business_id TEXT NOT NULL UNIQUE REFERENCES businesses(id) ON DELETE CASCADE,
  free_bonus INT NOT NULL DEFAULT 100,
  sub_credits INT NOT NULL DEFAULT 250,
  sub_price INT NOT NULL DEFAULT 250000,
  topup_price INT NOT NULL DEFAULT 1000,
  tier1_price INT NOT NULL DEFAULT 250000,
  tier1_credits INT NOT NULL DEFAULT 300,
  tier2_price INT NOT NULL DEFAULT 500000,
  tier2_credits INT NOT NULL DEFAULT 625,
  tier3_price INT NOT NULL DEFAULT 1000000,
  tier3_credits INT NOT NULL DEFAULT 1500,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
