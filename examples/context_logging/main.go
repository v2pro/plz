package main

import (
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/concurrent"
	"time"
)

func main() {
	concurrent.GlobalUnboundedExecutor.Go(func(ctx *countlog.Context) {
		ctx.Add("traceId", "axkenfppkl")
		processRequest(ctx, "request 111")
	})
	time.Sleep(time.Second)
}

func processRequest(ctx *countlog.Context, request string) {
	ctx.Add("userId", "111")
	ctx.Info("calculated game scores", "score", []int{1, 2, 3})
}
