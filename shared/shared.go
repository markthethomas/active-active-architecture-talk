package shared

import (
	"os"
)

var ENV string

func init() {
	ENV = os.Getenv("GO_ENV")
	InitLogging()
}
