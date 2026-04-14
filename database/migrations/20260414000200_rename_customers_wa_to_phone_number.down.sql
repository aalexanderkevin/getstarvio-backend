DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'customers' AND column_name = 'phone_number'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'customers' AND column_name = 'wa'
  ) THEN
    ALTER TABLE customers RENAME COLUMN phone_number TO wa;
  END IF;
END $$;

DROP INDEX IF EXISTS idx_customers_phone_number;
CREATE INDEX IF NOT EXISTS idx_customers_wa ON customers(wa);
