package main

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	samplingConfig := &zap.SamplingConfig{
		Initial:    3, // Allow first 3 events
		Thereafter: 3, // Allows 1 per 3 thereafter
		Hook: func(e zapcore.Entry, d zapcore.SamplingDecision) {
			if d == zapcore.LogDropped {
				fmt.Println("event dropped...")
			}
		},
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Sampling = samplingConfig
	cfg.EncoderConfig.TimeKey = "" // Turn off timestamp output

	logger, _ := cfg.Build()

	zap.ReplaceGlobals(logger) // Replace Zap's global logger
}

func main() {
	for i := 1; i <= 9; i++ {
		zap.S().Infow(
			"Testing sampling",
			"index", i,
		)
	}
}
