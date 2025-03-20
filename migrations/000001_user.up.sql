CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- users: хранение данных о пользователях
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       login TEXT UNIQUE NOT NULL,
                       pass_hash TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT NOW()
);