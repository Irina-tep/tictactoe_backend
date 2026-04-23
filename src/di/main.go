package main

import (
	"log"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {

	app := fx.New(
		fx.Provide(
			NewStorage,
			NewGameRepo,
			NewUserRepo,
			NewGameService,
			NewUserService,
			NewRouter,
			NewHTTPServer,
		),
		fx.Invoke(func(lc fx.Lifecycle, server *http.Server) {
			// Просто инициализация, сервер уже запущен в OnStart
		}),
		// Включаем логирование fx
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ConsoleLogger{W: log.Writer()}
		}),
	)
	app.Run()
}
