CREATE TABLE IF NOT EXISTS facebook_logs (
  id TEXT PRIMARY KEY,
  operation TEXT NOT NULL,
  url TEXT NOT NULL,
  request_body TEXT NOT NULL,
  response_body TEXT,
  response_code INT NOT NULL,
  ref_id TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_facebook_logs_operation ON facebook_logs(operation);
CREATE INDEX IF NOT EXISTS idx_facebook_logs_ref_id ON facebook_logs(ref_id);
CREATE INDEX IF NOT EXISTS idx_facebook_logs_created_at ON facebook_logs(created_at);
