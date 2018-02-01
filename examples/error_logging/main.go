package main

import (
	"context"
	"fmt"
	"github.com/v2pro/plz"
	"github.com/v2pro/plz/concurrent"
	"github.com/v2pro/plz/countlog"
	"os"
	"time"
	"github.com/v2pro/plz/countlog/output/compact"
)

func main() {
	executor := concurrent.GlobalUnboundedExecutor
	defer executor.StopAndWaitForever()
	countlog.Configure(func(cfg *countlog.Config) {
		cfg.EventWriterConfig.Executor = executor.Adapt()
		cfg.Format = &compact.Format{}
		cfg.Writer = os.Stdout
	})
	plz.PlugAndPlay()
	ctx := context.WithValue(context.Background(), "traceId", "abcd")
	//err := doSomething(ctx)
	//countlog.TraceCall("callee!main.doSomething", err, "ctx", ctx)
	doZ(countlog.Ctx(ctx))
}

func doX(ctx context.Context) error {
	file, err := os.OpenFile("/tmp/my-dir/abc", os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte("hello"))
	if err != nil {
		return err
	}
	return nil
}

func doA(ctx context.Context) error {
	file, err := os.OpenFile("/tmp/my-dir/abc", os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	_, err = file.Write([]byte("hello"))
	if err != nil {
		return fmt.Errorf("failed to write: %v", err)
	}
	return nil
}

func doY(ctx context.Context) error {
	defer func() {
		recovered := recover()
		if recovered != nil {
			countlog.Fatal("event!doY.panic",
				"err", recovered,
				"stacktrace", countlog.ProvideStacktrace)
		}
	}()
	start := time.Now()
	file, err := os.OpenFile("/tmp/my-dir/abc", os.O_RDWR, 0666)
	if err != nil {
		countlog.Error("event!metric",
			"callee", "ioutil.WriteFile", "ctx", ctx, "latency", time.Since(start))
		countlog.Error("event!doY.failed to open file", "err", err)
		return err
	}
	countlog.Trace("event!metric",
		"callee", "ioutil.WriteFile", "ctx", ctx, "latency", time.Since(start))
	defer func() {
		err = file.Close()
		if err != nil {
			countlog.Error("event!doY.failed to close file", "err", err)
		}
	}()
	_, err = file.Write([]byte("hello"))
	if err != nil {
		return err
	}
	return nil
}

func doZ(ctx *countlog.Context) error {
	defer func() {
		countlog.LogPanic(recover())
	}()
	file, err := os.OpenFile("/tmp/abc", os.O_RDWR, 0666)
	ctx.TraceCall("callee!os.OpenFile", err)
	if err != nil {
		return err
	}
	defer plz.Close(file)
	_, err = file.Write([]byte("hello"))
	ctx.TraceCall("callee!file.Write", err)
	if err != nil {
		return err
	}
	return nil
}
