package logger

import (
	"log"

	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

var (
	gelfLogger *gelf.UDPWriter
)

func InitGelfLogger(endpoint string) {
	var err error
	gelfLogger, err = gelf.NewUDPWriter(endpoint)
	if err != nil || gelfLogger == nil {
		log.Fatal("Graylog недоступен:", err)
	}
	gelfLogger.Facility = "auth-service"

}

func WriteMessage(message *gelf.Message) error {
	return gelfLogger.WriteMessage(message)
}

func Close() {
	gelfLogger.Close()
}
