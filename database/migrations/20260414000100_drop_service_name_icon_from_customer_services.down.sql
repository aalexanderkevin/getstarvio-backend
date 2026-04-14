ALTER TABLE customer_services
DROP CONSTRAINT IF EXISTS uq_customer_services_customer_category;

ALTER TABLE customer_services
ADD COLUMN IF NOT EXISTS service_name TEXT;

ALTER TABLE customer_services
ADD COLUMN IF NOT EXISTS service_icon TEXT;

UPDATE customer_services cs
SET
  service_name = COALESCE(cat.name, 'Layanan'),
  service_icon = COALESCE(cat.icon, '✨')
FROM categories cat
WHERE cat.id = cs.category_id;

UPDATE customer_services
SET
  service_name = COALESCE(service_name, 'Layanan'),
  service_icon = COALESCE(service_icon, '✨');

ALTER TABLE customer_services
ALTER COLUMN service_name SET NOT NULL;

ALTER TABLE customer_services
ALTER COLUMN service_icon SET NOT NULL;

ALTER TABLE customer_services
ADD CONSTRAINT customer_services_customer_id_service_name_key UNIQUE (customer_id, service_name);
