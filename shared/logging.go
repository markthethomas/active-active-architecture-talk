package shared

import (
	"go.uber.org/zap"
)

// loggingInitialized is a flag to keep track of whether logging has been set up
var loggingInitialized = false

// Logger is the global logger to use
var Logger zap.SugaredLogger

// InitLogging initializes the logger
func InitLogging() {
	if loggingInitialized == true {
		return
	}
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	Logger = *sugar
	loggingInitialized = true
}
