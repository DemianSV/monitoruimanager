CREATE TABLE IF NOT EXISTS auth_refresh_sessions (
    token_hash text PRIMARY KEY,
    user_name text,
    expires_at timestamp,
    consumed boolean
);
