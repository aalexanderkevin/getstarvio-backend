CREATE TABLE IF NOT EXISTS wa_templates (
  id TEXT PRIMARY KEY,
  meta_template_name TEXT NOT NULL,
  template_alias TEXT NOT NULL,
  category TEXT NOT NULL CHECK (category IN ('UTILITY', 'MARKETING', 'AUTHENTICATION')),
  language TEXT NOT NULL CHECK (language IN ('id', 'en_US', 'ms_MY')),
  status TEXT NOT NULL CHECK (status IN ('DRAFT', 'PENDING', 'APPROVED', 'REJECTED', 'PAUSED', 'FLAGGED')),
  body TEXT NOT NULL,
  body_example TEXT NOT NULL,
  meta_template_id TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wa_templates_status ON wa_templates(status);
CREATE INDEX IF NOT EXISTS idx_wa_templates_meta_template_id ON wa_templates(meta_template_id);
