package output

import (
	"context"
	"io"
	"github.com/v2pro/plz/countlog/spi"
	"fmt"
)

type Executor func(func(ctx context.Context))

func DefaultExecutor(handler func(ctx context.Context)) {
	go func() {
		handler(context.Background())
	}()
}

type blockingQueueWriter struct {
	queue  chan []byte
	writer io.Writer
}

type nonBlockingQueueWriter struct {
	blockingQueueWriter
	onMessageDropped func(message []byte)
}

type AsyncWriterConfig struct {
	Executor         Executor
	QueueLength      int
	Writer           io.Writer
	IsQueueBlocking  bool
	OnMessageDropped func(msg []byte)
}

func NewAsyncWriter(cfg AsyncWriterConfig) io.Writer {
	executor := cfg.Executor
	if executor == nil {
		executor = DefaultExecutor
	}
	queueLength := cfg.QueueLength
	if queueLength == 0 {
		queueLength = 1024
	}
	onMessageDropped := cfg.OnMessageDropped
	if onMessageDropped == nil {
		droppedCount := 0
		onMessageDropped = func(msg []byte) {
			droppedCount++
			if droppedCount%1000 == 1 {
				spi.OnError(fmt.Errorf("countlog async writer congestion, dropped %v messages so far", droppedCount))
			}
		}
	}
	if cfg.IsQueueBlocking {
		asyncWriter := &blockingQueueWriter{
			queue:  make(chan []byte, queueLength),
			writer: cfg.Writer,
		}
		executor(asyncWriter.asyncWrite)
		return asyncWriter
	}
	asyncWriter := &nonBlockingQueueWriter{
		blockingQueueWriter: blockingQueueWriter{
			queue:  make(chan []byte, queueLength),
			writer: cfg.Writer,
		},
		onMessageDropped: onMessageDropped,
	}
	return asyncWriter
}

func (writer *blockingQueueWriter) asyncWrite(ctx context.Context) {
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

func (writer *blockingQueueWriter) Write(buf []byte) (int, error) {
	writer.queue <- buf
	return len(buf), nil
}

func (writer *nonBlockingQueueWriter) Write(buf []byte) (int, error) {
	select {
	case writer.queue <- buf:
	default:
		writer.onMessageDropped(buf)
	}
	return len(buf), nil

}
