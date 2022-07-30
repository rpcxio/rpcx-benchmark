package main

import (
	"flag"
	"runtime"
	"time"

	"github.com/rpcxio/rpcx-benchmark/proto"
	rlog "github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"
)

func hello(ctx *server.Context) error {
	msg := &proto.BenchmarkMessage{}

	err := ctx.Bind(msg)
	if err != nil {
		return err
	}

	msg.Field1 = "OK"
	msg.Field2 = 100

	if *delay > 0 {
		time.Sleep(*delay)
	} else {
		runtime.Gosched()
	}

	return ctx.Write(msg)
}

var (
	host  = flag.String("s", "127.0.0.1:8972", "listened ip and port")
	delay = flag.Duration("delay", 0, "delay to mock business processing by sleep")
)

func main() {
	flag.Parse()

	rlog.SetDummyLogger()

	rpcxserver := server.NewServer(server.WithPool(100, 1000000))
	// rpcxserver.AsyncWrite = true

	rpcxserver.AddHandler("Hello", "Say", hello)
	rpcxserver.Serve("tcp", *host)
}
