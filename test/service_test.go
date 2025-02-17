package test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	exchange "github.com/AndrewTarev/proto-repo/gen/exchange"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"gw-exchanger/internal/delivery/grpc_delivery"
	"gw-exchanger/internal/service"
	"gw-exchanger/internal/service/mocks"
	"gw-exchanger/internal/storage/models"
)

func TestGetExchangeRates(t *testing.T) {
	// Создаем контроллер для gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для ExchangeService
	mockService := mocks.NewMockExchangeService(ctrl)

	// Настройка мок-сервиса для возвращения данных с decimal.Decimal
	mockService.EXPECT().GetExchangeRates(gomock.Any(), "RUB").Return(models.ExchangeRateResponse{
		Rates: map[string]decimal.Decimal{
			"USD": decimal.NewFromFloat(75.0),
			"EUR": decimal.NewFromFloat(90.0),
		},
	}, nil).Times(1)

	// Создаем экземпляр сервиса с мок-сервисом
	services := service.Service{ExchangeService: mockService}

	// Создаем хэндлер, передавая сервис
	handler := grpc_delivery.NewExchangerHandler(&services)

	// Создаем gRPC запрос
	req := &exchange.Empty{}

	// Вызываем метод хэндлера
	resp, err := handler.GetExchangeRates(context.Background(), req)

	// Проверяем на отсутствие ошибок
	assert.NoError(t, err)

	usdRate, err := strconv.ParseFloat(resp.Rates["USD"], 64)
	assert.NoError(t, err)
	assert.Equal(t, 75.0, usdRate)

	eurRate, err := strconv.ParseFloat(resp.Rates["EUR"], 64)
	assert.NoError(t, err)
	assert.Equal(t, 90.0, eurRate)
}

func TestGetExchangeRateForCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockExchangeService(ctrl)
	services := service.Service{ExchangeService: mockService}
	handler := grpc_delivery.NewExchangerHandler(&services)

	// Успешный запрос, валюта найдена
	mockService.EXPECT().
		GetExchangeRates(gomock.Any(), "USD").
		Return(models.ExchangeRateResponse{
			Rates: map[string]decimal.Decimal{
				"RUB": decimal.NewFromFloat(75.0),
				"EUR": decimal.NewFromFloat(90.0),
			},
		}, nil).
		Times(1)

	// Создаем gRPC запрос и ответ
	req := &exchange.CurrencyRequest{
		FromCurrency: "USD",
		ToCurrency:   "RUB",
	}
	resp, err := handler.GetExchangeRateForCurrency(context.Background(), req)

	// Проверки успешного выполнения
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "USD", resp.FromCurrency)
	assert.Equal(t, "RUB", resp.ToCurrency)

	// Используем strconv для парсинга строки в float64
	usdRate, err := strconv.ParseFloat(resp.Rate, 64)
	assert.NoError(t, err)
	assert.Equal(t, 75.0, usdRate) // Проверяем, что курс USD равен 75.0

	// Тестируем ошибку, когда валюта не найдена
	mockService.EXPECT().
		GetExchangeRates(gomock.Any(), "XYZ").
		Return(models.ExchangeRateResponse{}, errors.New("currency not found")).
		Times(1)

	req = &exchange.CurrencyRequest{
		FromCurrency: "XYZ",
		ToCurrency:   "USD",
	}
	resp, err = handler.GetExchangeRateForCurrency(context.Background(), req)

	// Проверки ошибки, когда валюта не найдена
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "currency not found")
	assert.Nil(t, resp)

	// Тестируем внутреннюю ошибку сервиса
	mockService.EXPECT().
		GetExchangeRates(gomock.Any(), "EUR").
		Return(models.ExchangeRateResponse{}, errors.New("internal error")).
		Times(1)

	req = &exchange.CurrencyRequest{
		FromCurrency: "EUR",
		ToCurrency:   "USD",
	}
	resp, err = handler.GetExchangeRateForCurrency(context.Background(), req)

	// Проверки внутренней ошибки
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "internal error")
	assert.Nil(t, resp)
}
