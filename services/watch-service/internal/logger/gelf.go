package logger

import (
	"log"

	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

var (
	gelfLogger *gelf.UDPWriter
)

func InitGelfLogger(endpoint string) {
	gelfLogger, err := gelf.NewUDPWriter(endpoint)
	gelfLogger.Facility = "watch-service"

	if err != nil {
		log.Fatal("Graylog недоступен:", err)
	}
}

func WriteMessage(message *gelf.Message) error {
	return gelfLogger.WriteMessage(message)
}

func Close() {
	gelfLogger.Close()
}
