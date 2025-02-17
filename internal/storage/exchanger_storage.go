package storage

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"gw-exchanger/internal/errs"
	"gw-exchanger/internal/storage/models"
)

type Exchanger struct {
	db *pgxpool.Pool
}

func NewExchanger(db *pgxpool.Pool) *Exchanger {
	return &Exchanger{
		db: db,
	}
}

// Поддерживаемые валюты
var supportedCurrencies = map[string]struct{}{
	"USD": {},
	"RUB": {},
	"EUR": {},
}

func (e *Exchanger) GetExchangeRates(ctx context.Context, baseCurrency string) (models.ExchangeRateResponse, error) {
	// Проверяем, поддерживается ли запрашиваемая валюта
	baseCurrency = strings.ToUpper(baseCurrency)
	if _, ok := supportedCurrencies[baseCurrency]; !ok {
		return models.ExchangeRateResponse{}, errs.ErrUnsupportedInputCurr
	}

	query := `
        SELECT base_currency, rate_usd, rate_rub, rate_eur 
        FROM exchange_rates 
        WHERE base_currency = $1 
        LIMIT 1;
    `

	var rateUSD, rateRUB, rateEUR decimal.Decimal
	err := e.db.QueryRow(ctx, query, baseCurrency).Scan(&baseCurrency, &rateUSD, &rateRUB, &rateEUR)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ExchangeRateResponse{}, errs.ErrNoRows
		}
		return models.ExchangeRateResponse{}, err
	}

	// Формируем ответ
	rates := make(map[string]decimal.Decimal)
	for currency := range supportedCurrencies {
		if currency == baseCurrency {
			continue
		}
		switch currency {
		case "USD":
			rates["USD"] = rateUSD
		case "RUB":
			rates["RUB"] = rateRUB
		case "EUR":
			rates["EUR"] = rateEUR
		}
	}

	return models.ExchangeRateResponse{
		Rates: rates,
	}, nil
}
