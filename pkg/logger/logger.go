package logger

import (
	"go.uber.org/zap"

	appErrors "github.com/PoorMercymain/user-segmenter/errors"
)

var instance *zap.SugaredLogger

func InitLogger() error {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"logfile.log", "stdout"}

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	instance = logger.Sugar()
	return nil
}

func GetLogger() (*zap.SugaredLogger, error) {
	var err error

	if instance == nil {
		err = appErrors.ErrorLoggerNotInitialized
	}

	return instance, err
}
