package output

import (
	"context"
	"io"
	"github.com/v2pro/plz/countlog/spi"
)

type Executor func(func(ctx context.Context))

func DefaultExecutor(handler func(ctx context.Context)) {
	go func() {
		handler(context.Background())
	}()
}

type asyncWriter struct {
	queue  chan []byte
	writer io.Writer
}

func newAsyncWriter(executor Executor, writer io.Writer) *asyncWriter {
	asyncWriter := &asyncWriter{
		queue:  make(chan []byte, 1024),
		writer: writer,
	}
	executor(asyncWriter.asyncWrite)
	return asyncWriter
}

func (writer *asyncWriter) asyncWrite(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case buf := <-writer.queue:
			_, err := writer.writer.Write(buf)
			if err != nil {
				spi.OnError(err)
			}
		}
	}
}

func (writer *asyncWriter) Write(buf []byte) (int, error) {
	select {
	case writer.queue <- buf:
	default:
		// TODO: handle queue overflow
	}
	return len(buf), nil
}
