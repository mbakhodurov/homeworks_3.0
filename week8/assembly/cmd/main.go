package main

import (
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/mbakhodurov/homeworks2/week8/assembly/internal/app"
	"github.com/mbakhodurov/homeworks2/week8/assembly/internal/config"
)

func main() {
	_ = godotenv.Load("assembly.env") //nolint:gosec // .env файл опционален — ошибка загрузки допустима

	config.MustLoad()

	if err := app.New().Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
	}
}
