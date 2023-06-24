package main

import (
	"fmt"
	"net/http"
	handler "notification/internal/controllers/handlers"
	log "notification/internal/platform/repositories"
	"notification/internal/usecase/notification"
)

var (
	url           = "../internal/logs.txt"
	logRepository *log.LogRepository
)

func main() {
	StartServer()

}

func StartServer() {
	logRepository = log.NewLogRepository(url)
	notificationApp := notification.NewNotificationUseCase(*logRepository)
	handler := handler.NewNotificationHandler(notificationApp)
	handler.RegisterRoutes()

	fmt.Println("Server listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
