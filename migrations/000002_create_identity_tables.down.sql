-- Drop identity tables in reverse order

DROP TABLE IF EXISTS identity.password_reset_tokens;
DROP TABLE IF EXISTS identity.refresh_tokens;
DROP TABLE IF EXISTS identity.users;
