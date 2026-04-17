INSERT INTO default_categories (
  id,
  name,
  icon,
  interval_days,
  template_id,
  template_body,
  example_body,
  is_active
)
VALUES
  ('defcat-facial-treatment', 'Facial Treatment', '💆', 30, 'tpl-a', 'Halo {{1}}! Sudah {{2}} hari sejak {{3}} terakhir kamu di {{4}}. Yuk balik lagi — kami tunggu! 😊', '["Pelanggan","{{interval}}","{{service}}","{{business}}"]', TRUE),
  ('defcat-waxing', 'Waxing', '🪒', 14, 'tpl-b', 'Halo {{1}}! Sudah {{2}} hari sejak {{3}} terakhir kamu di {{4}}. Yuk balik lagi — kami tunggu! 😊', '["Pelanggan","{{interval}}","{{service}}","{{business}}"]', TRUE),
  ('defcat-manicure-pedicure', 'Manicure & Pedicure', '💅', 21, 'tpl-a', 'Halo {{1}}! Sudah {{2}} hari sejak {{3}} terakhir kamu di {{4}}. Yuk balik lagi — kami tunggu! 😊', '["Pelanggan","{{interval}}","{{service}}","{{business}}"]', TRUE),
  ('defcat-body-massage', 'Body Massage', '🧖', 21, 'tpl-c', 'Halo {{1}}! Sudah {{2}} hari sejak {{3}} terakhir kamu di {{4}}. Yuk balik lagi — kami tunggu! 😊', '["Pelanggan","{{interval}}","{{service}}","{{business}}"]', TRUE),
  ('defcat-hair-treatment', 'Hair Treatment', '💇', 45, 'tpl-d', 'Halo {{1}}! Sudah {{2}} hari sejak {{3}} terakhir kamu di {{4}}. Yuk balik lagi — kami tunggu! 😊', '["Pelanggan","{{interval}}","{{service}}","{{business}}"]', TRUE),
  ('defcat-lash-lift-tint', 'Lash Lift & Tint', '👁️', 42, 'tpl-b', 'Halo {{1}}! Sudah {{2}} hari sejak {{3}} terakhir kamu di {{4}}. Yuk balik lagi — kami tunggu! 😊', '["Pelanggan","{{interval}}","{{service}}","{{business}}"]', TRUE)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  icon = EXCLUDED.icon,
  interval_days = EXCLUDED.interval_days,
  template_id = EXCLUDED.template_id,
  template_body = EXCLUDED.template_body,
  example_body = EXCLUDED.example_body,
  is_active = EXCLUDED.is_active,
  updated_at = NOW();
