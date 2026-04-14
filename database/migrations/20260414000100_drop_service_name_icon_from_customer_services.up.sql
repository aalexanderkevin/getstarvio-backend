ALTER TABLE customer_services
DROP CONSTRAINT IF EXISTS customer_services_customer_id_service_name_key;

ALTER TABLE customer_services
DROP COLUMN IF EXISTS service_name;

ALTER TABLE customer_services
DROP COLUMN IF EXISTS service_icon;

ALTER TABLE customer_services
ADD CONSTRAINT uq_customer_services_customer_category UNIQUE (customer_id, category_id);
