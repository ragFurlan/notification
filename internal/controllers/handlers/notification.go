package notification_handler

import (
	"encoding/json"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
		Message:  requestBody.Message,
		Category: requestBody.Category,
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logs, err := h.NotificationUseCase.GetLogs()
	if err != nil {
		http.Error(w, "Failed to get logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *NotificationHandler) DeleteLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.NotificationUseCase.DeleteLogs()
	if err != nil {
		http.Error(w, "Failed on delete logs", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Logs deleted",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *NotificationHandler) RegisterRoutes() {
	http.HandleFunc("/add", h.SubmitNotification)
	http.HandleFunc("/get", h.GetLogs)
	http.HandleFunc("/delete", h.DeleteLogs)
}
