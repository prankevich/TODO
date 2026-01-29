package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:01:05"}
	logger := zerolog.New(output).With().Timestamp().Logger()
	return logger

}
