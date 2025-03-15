CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- users: хранение данных о пользователях
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       username TEXT UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT NOW()
);

-- secrets: хранение приватных данных пользователей
CREATE TABLE secrets (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                         type TEXT NOT NULL, -- тип данных ("password", "text", "binary", "card")
                         data BYTEA NOT NULL, -- зашифрованные данные
                         metadata JSONB, -- метаинформация
                         created_at TIMESTAMP DEFAULT NOW(),
                         updated_at TIMESTAMP DEFAULT NOW()
);

-- devices: хранение информации об авторизованных клиентах
CREATE TABLE devices (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                         device_info TEXT NOT NULL, -- описание устройства (например, "Windows CLI")
                         last_seen TIMESTAMP DEFAULT NOW()
);

-- audit_log: логирование действий пользователей
CREATE TABLE audit_log (
                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                           user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           action TEXT NOT NULL, -- действие ("create_secret", "delete_secret", "update_secret", "login")
                           timestamp TIMESTAMP DEFAULT NOW()
);