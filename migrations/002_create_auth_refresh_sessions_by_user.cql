CREATE TABLE IF NOT EXISTS auth_refresh_sessions_by_user (
    user_name text,
    token_hash text,
    PRIMARY KEY (user_name, token_hash)
);
