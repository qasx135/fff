package logger

import (
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

var (
	gelfLogger *gelf.UDPWriter
)

func InitGelfLogger(endpoint string) error {
	var err error
	gelfLogger, err = gelf.NewUDPWriter(endpoint)
	gelfLogger.Facility = "catalog-service"

	return err
}

func WriteMessage(message *gelf.Message) error {
	return gelfLogger.WriteMessage(message)
}

func Close() {
	gelfLogger.Close()
}
