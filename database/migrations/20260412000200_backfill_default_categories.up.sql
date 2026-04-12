CREATE TABLE IF NOT EXISTS default_categories (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  icon TEXT NOT NULL,
  interval_days INT NOT NULL,
  template_id TEXT NOT NULL,
  template_body TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_default_categories_is_active ON default_categories(is_active);

INSERT INTO default_categories (
  id,
  name,
  icon,
  interval_days,
  template_id,
  template_body,
  is_active
)
VALUES
  ('defcat-facial-treatment', 'Facial Treatment', '💆', 30, 'tpl-a', 'Hai [nama], sudah sebulan sejak Facial Treatment terakhir kamu di [bisnis]. Yuk booking lagi biar kulit tetap glowing! ✨', TRUE),
  ('defcat-waxing', 'Waxing', '🪒', 14, 'tpl-b', 'Hai [nama]! Sudah 2 minggu nih sejak waxing terakhir. Mau reschedule? Hubungi kami ya 😊', TRUE),
  ('defcat-manicure-pedicure', 'Manicure & Pedicure', '💅', 21, 'tpl-a', 'Hai [nama], sudah 3 minggu sejak Manicure & Pedicure terakhir kamu di [bisnis]. Yuk booking lagi! 💅', TRUE),
  ('defcat-body-massage', 'Body Massage', '🧖', 21, 'tpl-c', 'Hai [nama], badan pegal? Sudah waktunya Body Massage lagi di [bisnis]. Ada diskon 10% kalau booking minggu ini! 💆‍♀️', TRUE),
  ('defcat-hair-treatment', 'Hair Treatment', '💇', 45, 'tpl-d', 'Hai [nama], gimana rambut kamu setelah Hair Treatment di [bisnis]? Kalau mau touch-up, kabari kami ya! 💇‍♀️', TRUE),
  ('defcat-lash-lift-tint', 'Lash Lift & Tint', '👁️', 42, 'tpl-b', 'Hai [nama]! Sudah lama nih sejak Lash Lift & Tint terakhir. Bulu mata kamu pasti kangen perawatan 😍', TRUE)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  icon = EXCLUDED.icon,
  interval_days = EXCLUDED.interval_days,
  template_id = EXCLUDED.template_id,
  template_body = EXCLUDED.template_body,
  is_active = EXCLUDED.is_active,
  updated_at = NOW();
