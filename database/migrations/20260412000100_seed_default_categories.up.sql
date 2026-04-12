WITH default_categories(name, icon, interval_days, template_id, template_body) AS (
  VALUES
    ('Facial Treatment', '💆', 30, 'tpl-a', 'Hai [nama], sudah sebulan sejak Facial Treatment terakhir kamu di [bisnis]. Yuk booking lagi biar kulit tetap glowing! ✨'),
    ('Waxing', '🪒', 14, 'tpl-b', 'Hai [nama]! Sudah 2 minggu nih sejak waxing terakhir. Mau reschedule? Hubungi kami ya 😊'),
    ('Manicure & Pedicure', '💅', 21, 'tpl-a', 'Hai [nama], sudah 3 minggu sejak Manicure & Pedicure terakhir kamu di [bisnis]. Yuk booking lagi! 💅'),
    ('Body Massage', '🧖', 21, 'tpl-c', 'Hai [nama], badan pegal? Sudah waktunya Body Massage lagi di [bisnis]. Ada diskon 10% kalau booking minggu ini! 💆‍♀️'),
    ('Hair Treatment', '💇', 45, 'tpl-d', 'Hai [nama], gimana rambut kamu setelah Hair Treatment di [bisnis]? Kalau mau touch-up, kabari kami ya! 💇‍♀️'),
    ('Lash Lift & Tint', '👁️', 42, 'tpl-b', 'Hai [nama]! Sudah lama nih sejak Lash Lift & Tint terakhir. Bulu mata kamu pasti kangen perawatan 😍')
),
target_businesses AS (
  SELECT b.id AS business_id
  FROM businesses b
  WHERE NOT EXISTS (
    SELECT 1
    FROM categories c
    WHERE c.business_id = b.id
  )
)
INSERT INTO categories (
  id,
  business_id,
  name,
  icon,
  interval_days,
  template_id,
  template_body,
  is_enabled
)
SELECT
  'seed-cat-' || substr(md5(tb.business_id || ':' || dc.name), 1, 24) AS id,
  tb.business_id,
  dc.name,
  dc.icon,
  dc.interval_days,
  dc.template_id,
  dc.template_body,
  TRUE AS is_enabled
FROM target_businesses tb
CROSS JOIN default_categories dc
ON CONFLICT (id) DO NOTHING;
