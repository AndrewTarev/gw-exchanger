# gw-exchanger

Cервис exchanger для получения курсов валют (gw-exchanger).
Сервис по grpc обрабатывает запросы на получение курса валют. В качестве валют поддерживаются USD, RUB, EUR.

Основное приложение для взаимодействия с сервисом - `gw-currency-wallet`:
- https://github.com/AndrewTarev/gw-currency-wallet

Поддерживаемые методы:
```
// Получение курсов обмена всех валют
rpc GetExchangeRates(Empty) returns (ExchangeRatesResponse);

// Получение курса обмена для конкретной валюты
rpc GetExchangeRateForCurrency(CurrencyRequest) returns (ExchangeRateResponse);
```

## Структура проекта
```
.
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── main.go
├── docker-compose.dev.yaml
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── app                                     // Сборка основных компонентов приложения
│   │   └── app.go
│   ├── config                                  // Конфигурационные файлы
│   │   ├── config.go
│   │   └── config.yaml                         // Дефолтный конфиг
│   ├── delivery                                // Слой работы с хэндлерами
│   │   └── grpc_delivery
│   │       ├── handler.go
│   │       └── middleware                      
│   │           └── err_middleware.go           // Обработка, логгирование ошибок и паник
│   ├── errs
│   │   └── errs_exchanger.go
│   ├── server                                  // Настройки grpc сервера
│   │   └── grpc.go
│   ├── service                                 // Слой сервисов 
│   │   ├── exchanger_service.go
│   │   ├── mocks
│   │   │   └── mock_service.go
│   │   └── service.go
│   └── storage                                 // Слой работы с БД
│       ├── exchanger_storage.go
│       ├── models                              // Сущноссти
│       │   └── exchange.go
│       └── storage.go
├── migrations                                  // Файлы миграций (миграции применяются через docker-compose)
│   ├── 000001_init.down.sql
│   └── 000001_init.up.sql
├── pkg                                         
│   ├── db                                      // Настройки Postgres
│   │   └── db.go
│   └── logging                                 // Настройки Logrus
│       └── logger.go
└── test                                        // Тесты
    └── service_test.go
```

## Установка приложения:

1. Склонируйте репозиторий себе на компьютер
    - git clone https://github.com/AndrewTarev/gw-exchanger.git
2. Установите свои переменные в .env файл
3. Запустите сборку контейнеров
    - docker-compose up --build