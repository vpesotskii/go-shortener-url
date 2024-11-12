package config

import (
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func Initialize(level string) error {
	atomicLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	config := zap.NewProductionConfig()
	config.Level = atomicLevel

	logger, err := config.Build()
	if err != nil {
		return err
	}

	Log = logger
	return nil
}
