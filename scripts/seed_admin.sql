-- Seed admin users (run this after migrations)
-- Replace phone numbers with actual admin phone numbers
INSERT INTO users (phone, name, role)
VALUES
    ('7979903861', 'Admin 1', 'admin'),
    ('9693392509', 'Admin 2', 'admin')
ON CONFLICT (phone) DO NOTHING;
