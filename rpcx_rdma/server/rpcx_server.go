package main

import (
	"context"
	"flag"
	"runtime"
	"time"

	"github.com/rpcxio/rpcx-benchmark/proto"
	rlog "github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"
)

type Hello int

func (t *Hello) Say(ctx context.Context, args *proto.BenchmarkMessage, reply *proto.BenchmarkMessage) error {
	args.Field1 = "OK"
	args.Field2 = 100
	*reply = *args
	if *delay > 0 {
		time.Sleep(*delay)
	} else {
		runtime.Gosched()
	}
	return nil
}

var (
	host  = flag.String("s", "127.0.0.1:8972", "listened ip and port")
	delay = flag.Duration("delay", 0, "delay to mock business processing by sleep")
)

func main() {
	flag.Parse()

	rlog.SetDummyLogger()

	rpcxserver := server.NewServer()
	rpcxserver.AsyncWrite = true

	rpcxserver.RegisterName("Hello", new(Hello), "")
	err := rpcxserver.Serve("rdma", *host)
	if err != nil {
		panic(err)
	}
}
