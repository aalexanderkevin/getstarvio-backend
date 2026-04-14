DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'customers' AND column_name = 'wa'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'customers' AND column_name = 'phone_number'
  ) THEN
    ALTER TABLE customers RENAME COLUMN wa TO phone_number;
  END IF;
END $$;

DROP INDEX IF EXISTS idx_customers_wa;
CREATE INDEX IF NOT EXISTS idx_customers_phone_number ON customers(phone_number);
