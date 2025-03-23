-- Удаление триггера
DROP TRIGGER IF EXISTS update_user_data_updated_at ON user_data;

-- Удаление функции
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаление индексов
DROP INDEX IF EXISTS idx_user_data_type;
DROP INDEX IF EXISTS idx_user_data_label;
DROP INDEX IF EXISTS idx_user_data_user_id;

-- Удаление таблицы
DROP TABLE IF EXISTS user_data;