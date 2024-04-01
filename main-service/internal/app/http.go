package app

import (
	"context"
	"go.uber.org/zap"
)

var prefix = "/api/v1"

func (app *App) StartHTTPServer(ctx context.Context) error {
	app.fiber.Static("/static", "./static")
	//grp := app.fiber.Group(prefix)
	// Создаем канал для сигнализации о завершении работы сервера
	done := make(chan struct{})
	// Запускаем сервер в горутине
	go func() {
		if err := app.fiber.Listen(app.Config.Port); err != nil {
			app.Logger.Fatal("error func StartHTTPServer, method Listen by path internal/app/http.go", zap.Error(err))
		}
		close(done) // Закрываем канал, когда сервер завершит работу
	}()
	select {
	case <-ctx.Done():
		// Получен сигнал о завершении работы, останавливаем сервер
		if err := app.fiber.Shutdown(); err != nil {
			app.Logger.Error("error shutting down the server", zap.Error(err))
		}
	case <-done:
		// Сервер завершил работу самостоятельно
		app.Logger.Info("server finished successfully")
	}
	return nil
}
