package notification_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notification/internal/entity"
	"notification/internal/usecase/notification"
)

type NotificationHandler struct {
	NotificationUseCase *notification.NotificationUseCase
}

func NewNotificationHandler(notificationUseCase *notification.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{
		NotificationUseCase: notificationUseCase,
	}
}

func (h *NotificationHandler) SubmitNotification(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Category entity.Category `json:"category"`
		Message  string          `json:"message"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	notification := entity.Notification{
		Message:   requestBody.Message,
		Category:  requestBody.Category,
		Notifiers: h.NotificationUseCase.GetNotifiers(requestBody.Category),
	}

	logs, err := h.NotificationUseCase.SendNotification(notification)
	if err != nil {
		http.Error(w, "Failed to send notification", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string       `json:"message"`
		Logs    []entity.Log `json:"logs"`
	}{
		Message: "Notification sent successfully",
		Logs:    logs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *NotificationHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	logs := h.NotificationUseCase.GetLogs()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *NotificationHandler) RegisterRoutes() {
	http.HandleFunc("/submit", h.SubmitNotification)
	http.HandleFunc("/logs", h.GetLogs)
}

func StartServer() {
	notificationApp := notification.NewNotificationUseCase()
	handler := NewNotificationHandler(notificationApp)
	handler.RegisterRoutes()

	fmt.Println("Server listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
