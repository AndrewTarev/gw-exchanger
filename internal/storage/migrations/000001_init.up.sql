CREATE TABLE exchange_rates (
    id SERIAL PRIMARY KEY,           -- Идентификатор записи
    base_currency VARCHAR(3) NOT NULL, -- Код базовой валюты (например, USD, EUR, RUB)
    rate_usd DECIMAL(10, 4) NOT NULL,  -- Курс по отношению к USD
    rate_rub DECIMAL(10, 4) NOT NULL,  -- Курс по отношению к RUB
    rate_eur DECIMAL(10, 4) NOT NULL,  -- Курс по отношению к EUR
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время создания записи
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Время последнего обновления записи
);

-- Индексы для ускорения запросов
CREATE INDEX idx_base_currency ON exchange_rates (base_currency);

INSERT INTO exchange_rates (base_currency, rate_usd, rate_rub, rate_eur)
VALUES
    ('USD', 1.0000, 75.5000, 0.9300),
    ('EUR', 1.0700, 80.5000, 1.0000),
    ('RUB', 0.0133, 1.0000, 0.0116);

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_updated_at
    BEFORE UPDATE ON exchange_rates
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();