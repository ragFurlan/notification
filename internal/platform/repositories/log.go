package log

import (
	"fmt"
	"notification/internal/entity"
)

type LogRepository struct{}

func (r *LogRepository) SaveLog(log entity.Log) {
	fmt.Printf("Saving log: %s\n", log.Message)
	// TODO: Logic to save the log to the database or log file
}

func (r *LogRepository) GetLogs() []entity.Log {
	return []entity.Log{{Message: "return log"}}
	// TODO:Logic to save the log to the database or log file
}
