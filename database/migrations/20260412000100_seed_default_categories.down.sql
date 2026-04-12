WITH default_categories(name) AS (
  VALUES
    ('Facial Treatment'),
    ('Waxing'),
    ('Manicure & Pedicure'),
    ('Body Massage'),
    ('Hair Treatment'),
    ('Lash Lift & Tint')
),
seeded_ids AS (
  SELECT
    'seed-cat-' || substr(md5(b.id || ':' || dc.name), 1, 24) AS id
  FROM businesses b
  CROSS JOIN default_categories dc
)
DELETE FROM categories c
USING seeded_ids s
WHERE c.id = s.id;
