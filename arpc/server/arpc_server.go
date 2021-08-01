package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"github.com/lesismal/arpc"
	"github.com/rpcxio/rpcx-benchmark/arpc/codec"
	"github.com/rpcxio/rpcx-benchmark/proto"
)

type Hello int

func (t *Hello) Say(ctx *arpc.Context) {
	args := &proto.BenchmarkMessage{}
	reply := &proto.BenchmarkMessage{}

	if err := ctx.Bind(args); err != nil {
		ctx.Error(err)
	}

	args.Field1 = "OK"
	args.Field2 = 100
	*reply = *args
	if *delay > 0 {
		time.Sleep(*delay)
	} else {
		runtime.Gosched()
	}
	ctx.Write(reply)
}

var (
	host       = flag.String("s", "127.0.0.1:8972", "listened ip and port")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	delay      = flag.Duration("delay", 0, "delay to mock business processing")
	debugAddr  = flag.String("d", "127.0.0.1:9981", "server ip and port")
	async      = flag.Bool("a", false, "async response flag")
)

func main() {
	flag.Parse()

	// alog.SetLogLevel(alog.LogLevelNone)

	log.Println("async response:", *async)

	go func() {
		log.Println(http.ListenAndServe(*debugAddr, nil))
	}()

	server := arpc.NewServer()
	server.Codec = &codec.ProtoBuffer{}
	server.Handler.Handle("Hello.Say", new(Hello).Say, *async)
	server.Run(*host)
}
