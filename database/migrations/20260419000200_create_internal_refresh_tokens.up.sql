CREATE TABLE IF NOT EXISTS internal_refresh_tokens (
  id TEXT PRIMARY KEY,
  admin_id TEXT NOT NULL REFERENCES internal_admins(id) ON DELETE CASCADE,
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_internal_refresh_tokens_admin_id ON internal_refresh_tokens(admin_id);
