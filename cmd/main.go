package main

import (
	"fmt"
	"net/http"
	handler "notification/internal/controllers/handlers"
	log "notification/internal/platform/repositories"
	"notification/internal/usecase/notification"
)

var (
	url                 = "../internal/logs.txt"
	notificationUseCase *notification.NotificationUseCase
)

func main() {
	logRepository := log.NewLogRepository(url)
	notificationUseCase = notification.NewNotificationUseCase(*logRepository)
	StartServer()

}

func StartServer() {
	handler := handler.NewNotificationHandler(notificationUseCase)
	handler.RegisterRoutes()

	fmt.Println("Server listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
