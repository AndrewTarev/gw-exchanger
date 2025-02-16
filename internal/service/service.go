package service

import (
	"context"

	"gw-exchanger/internal/storage"
	"gw-exchanger/internal/storage/models"
)

type ExchangeService interface {
	GetExchangeRates(ctx context.Context, baseCurrency string) (models.ExchangeRateResponse, error)
}

type Service struct {
	ExchangeService
}

func NewExchangerService(stor *storage.Storage) *Service {
	return &Service{
		ExchangeService: NewExchanger(stor),
	}
}
