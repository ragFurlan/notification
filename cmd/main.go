package main

import (
	"fmt"
	"net/http"
	controller "notification/internal/controllers/handlers"
	log "notification/internal/platform/repositories"
	"notification/internal/usecase/notification"

	"github.com/gorilla/handlers"
)

var (
	url                 = "../internal/logs.txt"
	notificationUseCase *notification.NotificationUseCase
)

func main() {
	logRepository := log.NewLogRepository(url)
	notificationUseCase = notification.NewNotificationUseCase(logRepository)
	StartServer()

}

func StartServer() {
	handler := controller.NewNotificationHandler(notificationUseCase)
	router := handler.RegisterRoutes()

	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	fmt.Println("Server listening on http://localhost:8080")
	http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(router))
}
