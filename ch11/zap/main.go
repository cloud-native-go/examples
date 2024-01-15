/*
 * Copyright 2024 Matthew A. Titmus
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
