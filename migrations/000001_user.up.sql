CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- users: хранение данных о пользователях
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       username TEXT UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT NOW()
);