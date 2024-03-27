package app

import (
	"go.uber.org/zap"
)

var prefix = "/api/v1"

func (app *App) StartHTTPServer() error {
	app.fiber.Static("/static", "./static")
	//grp := app.fiber.Group(prefix)
	if err := app.fiber.Listen(app.config.Port); err != nil {
		app.Logger.Fatal("error func StartHTTPServer, method Listen by path internal/app/http.go", zap.Error(err))
	}
	return nil
}
