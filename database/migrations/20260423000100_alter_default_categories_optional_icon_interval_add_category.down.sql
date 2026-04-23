UPDATE default_categories
SET icon = '✨'
WHERE icon IS NULL;

UPDATE default_categories
SET interval_days = 30
WHERE interval_days IS NULL;

ALTER TABLE default_categories
ALTER COLUMN icon SET NOT NULL;

ALTER TABLE default_categories
ALTER COLUMN interval_days SET NOT NULL;

ALTER TABLE default_categories
DROP COLUMN IF EXISTS category;
