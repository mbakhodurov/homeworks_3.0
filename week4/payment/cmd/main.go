package main

import (
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/mbakhodurov/homeworks2/week4/payment/internal/app"
	"github.com/mbakhodurov/homeworks2/week4/payment/internal/config"
)

func main() {
	_ = godotenv.Load("payment.env") //nolint:gosec // .env файл опционален — ошибка загрузки допустима

	config.MustLoad()

	if err := app.New().Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
	}
}
