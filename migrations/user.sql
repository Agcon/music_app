CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    user_role VARCHAR(20) NOT NULL DEFAULT 'user'
);

INSERT INTO users (username, email, password_hash, user_role)
VALUES ('admin', 'admin@admin.com', '$2a$10$UZV0aXqwFSD5.LjA9DOaCueozBK55OedLdQOvJ9wfSNJGB.JlTvPi', 'admin')
ON CONFLICT (email) DO NOTHING;