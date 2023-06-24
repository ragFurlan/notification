package log

import (
	"bufio"
	"fmt"
	"notification/internal/entity"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type LogRepository struct {
	logFilePath string
}

type Log interface {
	SaveLog(log entity.Log) error
	GetLogs() ([]entity.Log, error)
	DeleteLogs() error
}

func NewLogRepository(logFilePath string) Log {
	return &LogRepository{
		logFilePath: logFilePath,
	}
}

func (r *LogRepository) SaveLog(log entity.Log) error {
	file, err := os.OpenFile(r.logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open log file: %v", err)
	}
	defer file.Close()

	logEntry := fmt.Sprintf("Timestamp: %s|Category: %s|Notification Type: %s|Message: %s|ID: %s|UserID: %v\n",
		log.Timestamp.Format(time.RFC3339), log.Category, log.NotificationType, log.Message, log.ID, log.UserID)

	if _, err := file.WriteString(logEntry); err != nil {
		return fmt.Errorf("Failed to write log entry: %v", err)
	}

	return nil
}

func (r *LogRepository) GetLogs() ([]entity.Log, error) {
	file, err := os.Open(r.logFilePath)
	if err != nil {
		err := fmt.Errorf("Failed to open log file: %v", err)
		return nil, err
	}
	defer file.Close()

	var logs []entity.Log
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		logEntry := scanner.Text()
		log, err := parseLogEntry(logEntry)
		if err != nil {
			fmt.Printf("Failed to parse log entry: %v", err)
			continue
		}
		logs = append(logs, log)
	}

	if err := scanner.Err(); err != nil {
		err = fmt.Errorf("Failed to scan log file: %v", err)
		return nil, err
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Timestamp.After(logs[j].Timestamp)
	})

	return logs, nil
}

func (r *LogRepository) DeleteLogs() error {
	err := os.Remove(r.logFilePath)
	if err != nil {
		return err
	}
	return nil
}

func parseLogEntry(logEntry string) (entity.Log, error) {
	var log entity.Log

	// Split the log entry string by newline characters
	lines := strings.Split(logEntry, "|")

	// Iterate through each line and parse the log fields
	for _, line := range lines {
		// Extract the key-value pairs from each line
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			return log, fmt.Errorf("invalid log entry format: %s", logEntry)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set the log field based on the key
		switch key {
		case "Timestamp":
			timestamp, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return log, fmt.Errorf("failed to parse timestamp in log entry: %s", logEntry)
			}
			log.Timestamp = timestamp
		case "Category":
			log.Category = entity.Category(value)
		case "Notification Type":
			log.NotificationType = value
		case "Message":
			log.Message = value
		case "UserID":
			userID, _ := strconv.Atoi(value)
			log.UserID = userID
		case "ID":
			log.ID = value

		}
	}

	return log, nil
}
