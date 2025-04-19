CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
    );

INSERT INTO users (username, email, password_hash)
VALUES ('admin', 'admin@admin.com', '$2a$10$Q3Wa5kG0UEnA2c.6kspcu.3BMeqB0EHE2AvY38Zkq8nKDMI8EqjSe')
ON CONFLICT (email) DO NOTHING;