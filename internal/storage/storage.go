package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"gw-exchanger/internal/storage/models"
)

type ExchangerStorage interface {
	GetExchangeRates(ctx context.Context, baseCurrency string) (models.ExchangeRateResponse, error)
}

type Storage struct {
	ExchangerStorage
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{
		ExchangerStorage: NewExchanger(db),
	}
}
