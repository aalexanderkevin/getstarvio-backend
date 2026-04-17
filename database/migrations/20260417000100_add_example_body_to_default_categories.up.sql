ALTER TABLE default_categories
ADD COLUMN IF NOT EXISTS example_body TEXT;

UPDATE default_categories
SET example_body = '["Pelanggan","{{interval}}","{{service}}","{{business}}"]'
WHERE example_body IS NULL OR btrim(example_body) = '';

ALTER TABLE default_categories
ALTER COLUMN example_body SET NOT NULL;
