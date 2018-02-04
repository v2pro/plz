package concurrent

import (
	"github.com/v2pro/plz/countlog/output"
)

type ClosableWriter struct {
	output.ClosableWriter
	executor *UnboundedExecutor
}

func NewAsyncWriter(cfg output.AsyncWriterConfig) *ClosableWriter {
	executor := NewUnboundedExecutor()
	cfg.Executor = executor.Adapt()
	writer := output.NewAsyncWriter(cfg)
	return &ClosableWriter{
		ClosableWriter: writer,
		executor: executor,
	}
}

func (writer *ClosableWriter) Close() error {
	writer.executor.StopAndWaitForever()
	return writer.ClosableWriter.Close()
}