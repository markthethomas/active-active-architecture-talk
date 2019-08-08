package shared

import (
	"os"

	"gopkg.in/underarmour/dynago.v2"
)

// DAL is the global DB instance
var DAL *dynago.Client

var connected = false

// Connect connects to the DB and sets the global DB instance up
func Connect() {
	if connected {
		return
	}

	// connect to dynamo
	cfg := dynago.ExecutorConfig{
		Region:       os.Getenv("AWS_REGION"),
		AccessKey:    os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SessionToken: os.Getenv("AWS_SESSION_TOKEN"),
	}
	exec := dynago.NewAwsExecutor(cfg)
	c := dynago.NewClient(exec)
	DAL = c

	connected = true
}
