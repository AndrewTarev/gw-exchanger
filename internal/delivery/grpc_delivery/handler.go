package grpc_delivery

import (
	"context"
	"fmt"

	exch "github.com/AndrewTarev/proto-repo/gen/exchange"

	"gw-exchanger/internal/service"
)

type ExchangerHandler struct {
	exch.UnimplementedExchangeServiceServer
	service *service.Service
}

func NewExchangerHandler(service *service.Service) *ExchangerHandler {
	return &ExchangerHandler{
		service: service,
	}
}

// GetExchangeRates - получение всех курсов валют
func (h *ExchangerHandler) GetExchangeRates(ctx context.Context, req *exch.Empty) (*exch.ExchangeRatesResponse, error) {
	rates, err := h.service.ExchangeService.GetExchangeRates(ctx, "RUB")
	if err != nil {
		return nil, err
	}

	// Преобразуем в gRPC-ответ
	response := &exch.ExchangeRatesResponse{
		Rates: make(map[string]string),
	}

	for currency, rate := range rates.Rates {
		response.Rates[currency] = rate.String() // decimal.Decimal -> string
	}

	return response, nil
}

// GetExchangeRateForCurrency - получение курса одной валюты
func (h *ExchangerHandler) GetExchangeRateForCurrency(ctx context.Context, req *exch.CurrencyRequest) (*exch.ExchangeRateResponse, error) {
	rates, err := h.service.ExchangeService.GetExchangeRates(ctx, req.FromCurrency)
	if err != nil {
		return nil, err
	}

	// Находим нужный курс
	rate, ok := rates.Rates[req.ToCurrency]
	if !ok {
		return nil, fmt.Errorf("currency not found: %s", req.ToCurrency)
	}

	return &exch.ExchangeRateResponse{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         rate.String(),
	}, nil
}
