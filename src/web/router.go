package web

import (
	"datasource"
	"net/http"
)


func NewRouter(gameService *datasource.GameService, userService *datasource.AuthorizationService) http.Handler {
	handler := &GameHandler{gameService: gameService, userService: userService} //создает обработчик игры (GameHandler) с сервисом для работы с данными и с пользователем

	// Создаём middleware для проверки авторизации
	authMiddleware := handler.UserAuthenticator

	// Маршрутизатор http.ServeMux — это "диспетчер", который: принимает входящий HTTP-запрос (метод + URL), сравнивает его с зарегистрированными маршрутами, направляет запрос к нужной функции-обработчику
	mux := http.NewServeMux()

	
	mux.HandleFunc("POST /auth/signup", handler.SignUp)
	mux.HandleFunc("POST /auth/login", handler.Authenticate)
	
	mux.Handle("POST /game", authMiddleware(http.HandlerFunc(handler.CreateGame)))
	mux.Handle("GET /game/{id}", authMiddleware(http.HandlerFunc(handler.GetGame)))
	mux.Handle("GET /games", authMiddleware(http.HandlerFunc(handler.GetGames)))
	mux.Handle("GET /info/{id}", authMiddleware(http.HandlerFunc(handler.GetInfo)))
	mux.Handle("POST /game/{id}", authMiddleware(http.HandlerFunc(handler.UpdateGame)))
	mux.Handle("POST /game/{id}/join", authMiddleware(http.HandlerFunc(handler.JoinGame)))
	return mux //Возвращает готовый маршрутизатор, который можно использовать в HTTP-сервере
}
