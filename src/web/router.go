package web

import (
	"datasource"
	"net/http"
)

// создает маршрутизатор для API игры.
func NewRouter(gameService *datasource.GameService) http.Handler {
	handler := &GameHandler{gameService: gameService} //создает обработчик игры (GameHandler) с сервисом для работы с данными
	// Маршрутизатор http.ServeMux — это "диспетчер", который: принимает входящий HTTP-запрос (метод + URL), сравнивает его с зарегистрированными маршрутами, направляет запрос к нужной функции-обработчику
	mux := http.NewServeMux()
	//Настраивает маршруты — связывает URL-пути с конкретными функциями-обработчиками
	mux.HandleFunc("POST /game", handler.CreateGame)
	mux.HandleFunc("GET /game/{id}", handler.GetGame)
	mux.HandleFunc("POST /game/{id}", handler.UpdateGame)

	return mux //Возвращает готовый маршрутизатор, который можно использовать в HTTP-сервере
}
