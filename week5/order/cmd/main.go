package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/mbakhodurov/homeworks2/week5/order/internal/app"
	"github.com/mbakhodurov/homeworks2/week5/order/internal/config"
)

func main() {
	// .env опционален — ошибка загрузки допустима.
	// Загружаем до ResolveConfigPath/MustLoad, чтобы CONFIG_PATH из .env подхватился.
	_ = godotenv.Load("order.env") //nolint:gosec // .env файл опционален — ошибка загрузки допустима

	config.MustLoad()

	a := app.New(context.Background())

	if err := a.Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
	}
}
