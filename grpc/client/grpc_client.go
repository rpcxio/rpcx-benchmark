package main

import (
	"flag"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/juju/ratelimit"
	benchmark "github.com/rpcxio/rpcx-benchmark"
	"github.com/rpcxio/rpcx-benchmark/grpc/pb"
	"github.com/smallnest/rpcx/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var concurrency = flag.Int("c", 1, "concurrency")
var total = flag.Int("n", 10000, "total requests for all clients")
var host = flag.String("s", "127.0.0.1:8972", "server ip and port")
var pool = flag.Int("pool", 10, " shared grpc clients")
var rate = flag.Int("r", 10000, "throughputs")

func main() {
	flag.Parse()

	tb := ratelimit.NewBucket(time.Second/time.Duration(*rate), int64(*rate))

	// 并发goroutine数.模拟客户端
	n := *concurrency
	// 每个客户端需要发送的请求数
	m := *total / n
	log.Infof("concurrency: %d\nrequests per client: %d\n\n", n, m)

	servers := strings.Split(*host, ",")
	log.Infof("Servers: %+v\n\n", *host)

	args := prepareArgs()

	// 请求消息大小
	b, _ := proto.Marshal(args)
	log.Infof("message size: %d bytes\n\n", len(b))

	// 等待所有测试完成
	var wg sync.WaitGroup
	wg.Add(n * m)

	// 总请求数
	var trans uint64
	// 返回正常的总请求数
	var transOK uint64

	// 每个goroutine的耗时记录
	d := make([][]int64, n, n)

	// 创建客户端连接池
	var clientIndex uint64
	var poolClients = make([]pb.HelloClient, 0, *pool)
	for i := 0; i < *pool; i++ {
		conn, err := grpc.Dial(servers[0], grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		c := pb.NewHelloClient(conn)
		//warmup
		for j := 0; j < 5; j++ {
			c.Say(context.Background(), args)
		}
		poolClients = append(poolClients, c)
	}

	// 栅栏，控制客户端同时开始测试
	var startWg sync.WaitGroup
	startWg.Add(n + 1) // +1 是因为有一个goroutine用来记录开始时间

	// 创建客户端 goroutine 并进行测试
	var startTime = time.Now().UnixNano()
	go func() {
		startWg.Done()
		startWg.Wait()
		startTime = time.Now().UnixNano()
	}()
	for i := 0; i < n; i++ {
		dt := make([]int64, 0, m)
		d = append(d, dt)

		go func(i int) {
			for j := 0; j < m; j++ {
				// 限流，这里不把限流的时间计算到等待耗时中
				tb.Wait(1)

				t := time.Now().UnixNano()
				ci := atomic.AddUint64(&clientIndex, 1)
				ci = ci % uint64(*pool)
				c := poolClients[int(ci)]
				reply, err := c.Say(context.Background(), args)
				t = time.Now().UnixNano() - t // 等待时间+服务时间，等待时间是客户端调度的等待时间以及服务端读取请求、调度的时间，服务时间是请求被服务处理的实际时间

				d[i] = append(d[i], t)

				if err == nil && reply.Field1 == "OK" {
					atomic.AddUint64(&transOK, 1)
				}

				atomic.AddUint64(&trans, 1)
				wg.Done()
			}

		}(i)

	}

	wg.Wait()

	// 统计
	benchmark.Stats(startTime, *total, d, trans, transOK)
}

func prepareArgs() *pb.BenchmarkMessage {
	b := true
	var i int32 = 100000
	var i64 int64 = 100000
	var s = "许多往事在眼前一幕一幕，变的那麼模糊"

	var args pb.BenchmarkMessage

	v := reflect.ValueOf(&args).Elem()
	num := v.NumField()
	for k := 0; k < num; k++ {
		field := v.Field(k)
		if field.Type().Kind() == reflect.Ptr {
			switch v.Field(k).Type().Elem().Kind() {
			case reflect.Int, reflect.Int32:
				field.Set(reflect.ValueOf(&i))
			case reflect.Int64:
				field.Set(reflect.ValueOf(&i64))
			case reflect.Bool:
				field.Set(reflect.ValueOf(&b))
			case reflect.String:
				field.Set(reflect.ValueOf(&s))
			}
		} else {
			switch field.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64:
				field.SetInt(100000)
			case reflect.Bool:
				field.SetBool(true)
			case reflect.String:
				field.SetString(s)
			}
		}

	}
	return &args
}
