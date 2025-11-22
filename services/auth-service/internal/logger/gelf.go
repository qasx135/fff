package logger

import (
	"os"

	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

var (
	gelfLogger        *gelf.UDPWriter
	gelfWriterAddress = "graylog:12201"
)

func InitGelfLogger() error {
	if envGelfWriterAddress := os.Getenv("GELF_ENDPOINT"); envGelfWriterAddress != "" {
		gelfWriterAddress = envGelfWriterAddress
	}
	var err error
	gelfLogger, err = gelf.NewUDPWriter(gelfWriterAddress)
	gelfLogger.Facility = "auth-service"

	return err
}

func WriteMessage(message *gelf.Message) error {
	return gelfLogger.WriteMessage(message)
}
