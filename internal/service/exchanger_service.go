package service

import (
	"context"

	"gw-exchanger/internal/storage"
	"gw-exchanger/internal/storage/models"
)

type Exchanger struct {
	stor *storage.Storage
}

func NewExchanger(stor *storage.Storage) *Exchanger {
	return &Exchanger{
		stor: stor,
	}
}

func (e *Exchanger) GetExchangeRates(ctx context.Context, baseCurrency string) (models.ExchangeRateResponse, error) {
	if baseCurrency == "" {
		baseCurrency = "RUB"
	}

	resultCh := make(chan models.ExchangeRateResponse, 1)
	errCh := make(chan error, 1)

	go func() {
		rates, err := e.stor.GetExchangeRates(ctx, baseCurrency)
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- rates
	}()

	select {
	case result := <-resultCh:
		return result, nil
	case err := <-errCh:
		return models.ExchangeRateResponse{}, err
	case <-ctx.Done():
		return models.ExchangeRateResponse{}, ctx.Err()
	}
}
