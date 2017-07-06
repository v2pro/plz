package example

import "github.com/v2pro/plz"

var errorLogger = plz.LoggerOf("name", "error")

func Example_logging() {
	errorLogger.Error("shit happens", "user", 123)
	// Output:
}
