package countlog

import (
	"context"
	"io"
)

type Executor interface {
	Go(handler func(ctx *Context))
}

type defaultExecutor struct {
}

func (executor *defaultExecutor) Go(handler func(ctx *Context)) {
	go func() {
		handler(Ctx(context.Background()))
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
	executor.Go(asyncWriter.asyncWrite)
	return asyncWriter
}

func (writer *asyncWriter) asyncWrite(ctx *Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case buf := <-writer.queue:
			_, err := writer.writer.Write(buf)
			if err != nil {
				// TODO: handle error
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
