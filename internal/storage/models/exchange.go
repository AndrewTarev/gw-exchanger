package models

import "github.com/shopspring/decimal"

type ExchangeRateResponse struct {
	Rates map[string]decimal.Decimal `json:"rates"`
}

type ExchangeRatesResponse struct {
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
	Rate         string `json:"rates"`
}
