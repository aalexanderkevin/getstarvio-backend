ALTER TABLE facebook_logs
ADD COLUMN IF NOT EXISTS id_text TEXT;

UPDATE facebook_logs
SET id_text = id::TEXT
WHERE id_text IS NULL;

ALTER TABLE facebook_logs
DROP CONSTRAINT IF EXISTS facebook_logs_pkey;

ALTER TABLE facebook_logs
DROP COLUMN IF EXISTS id;

ALTER TABLE facebook_logs
RENAME COLUMN id_text TO id;

ALTER TABLE facebook_logs
ALTER COLUMN id SET NOT NULL;

ALTER TABLE facebook_logs
ADD PRIMARY KEY (id);
