-- Удаление триггера
DROP TRIGGER IF EXISTS trigger_update_updated_at ON exchange_rates;

-- Удаление функции
DROP FUNCTION IF EXISTS update_updated_at;

-- Удаление индекса
DROP INDEX IF EXISTS idx_base_currency;

-- Удаление данных из таблицы (по желанию)
DELETE FROM exchange_rates;

-- Удаление таблицы
DROP TABLE IF EXISTS exchange_rates;