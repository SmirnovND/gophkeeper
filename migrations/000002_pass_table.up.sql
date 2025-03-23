-- user_data: хранение паролей, данных карт и метаданных пользователей
CREATE TABLE user_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label TEXT NOT NULL,
    type TEXT NOT NULL, -- тип данных: 'password', 'card', 'metadata'
    data JSONB NOT NULL, -- данные в формате JSON
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (user_id, label) -- обеспечивает уникальность label в рамках одного пользователя
);

-- Индекс для ускорения поиска по user_id
CREATE INDEX idx_user_data_user_id ON user_data(user_id);

-- Индекс для ускорения поиска по label
CREATE INDEX idx_user_data_label ON user_data(label);

-- Индекс для ускорения поиска по типу данных
CREATE INDEX idx_user_data_type ON user_data(type);

-- Триггер для автоматического обновления updated_at при изменении записи
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_data_updated_at
    BEFORE UPDATE ON user_data
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();