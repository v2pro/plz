package countlog

import (
	"io"
	"runtime"
	"fmt"
)

func Close(closer io.Closer) {
	err := closer.Close()
	_, file, line, _ := runtime.Caller(1)
	TraceMetric("callee!closer", err, "closedAt", fmt.Sprintf("%s:%d", file, line))
}
