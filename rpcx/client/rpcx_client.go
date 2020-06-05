package main

import (
	"context"
	"flag"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/juju/ratelimit"
	"github.com/montanaflynn/stats"
	"github.com/rpcx-ecosystem/rpcx-benchmark/proto"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/protocol"
)

var concurrency = flag.Int("c", 1, "concurrency")
var total = flag.Int("n", 1, "total requests for all clients")
var host = flag.String("s", "127.0.0.1:8972", "server ip and port")
var pool = flag.Int("pool", 10, "shared rpcx clients")
var rate = flag.Int("r", 10000, "throughputs")

func main() {
	flag.Parse()

	tb := ratelimit.NewBucket(time.Second/time.Duration(*rate), int64(*rate))

	// 并发goroutine数.模拟客户端
	n := *concurrency
	// 每个客户端需要发送的请求数
	m := *total / n
	log.Infof("concurrency: %d\nrequests per client: %d\n\n", n, m)

	// 创建服务端的信息
	servers := strings.Split(*host, ",")
	var serverPeers []*client.KVPair
	for _, server := range servers {
		serverPeers = append(serverPeers, &client.KVPair{Key: server})
	}
	log.Infof("Servers: %+v\n\n", serverPeers)

	servicePath := "Hello"
	serviceMethod := "Say"

	// 准备好参数
	args := prepareArgs()

	// 参数的大小
	b := make([]byte, 1024)
	i, _ := args.MarshalTo(b)
	log.Infof("message size: %d bytes\n\n", i)

	var wg sync.WaitGroup
	wg.Add(n * m)

	// 创建客户端连接池
	var clientIndex uint64
	var poolClients = make([]client.XClient, 0, *pool)
	dis := client.NewMultipleServersDiscovery(serverPeers)
	for i := 0; i < *pool; i++ {
		option := client.DefaultOption
		option.SerializeType = protocol.ProtoBuffer
		xclient := client.NewXClient(servicePath, client.Failtry, client.RoundRobin, dis, option)
		defer xclient.Close()

		//warmup
		var reply proto.BenchmarkMessage
		for j := 0; j < 5; j++ {
			xclient.Call(context.Background(), serviceMethod, args, &reply)
		}

		poolClients = append(poolClients, xclient)
	}

	// 栅栏，控制客户端同时开始测试
	var startWg sync.WaitGroup
	startWg.Add(n)

	//  请求数
	var trans uint64
	// 正常的请求数
	var transOK uint64

	d := make([][]int64, n, n)

	// 创建客户端 goroutine 并进行测试
	totalT := time.Now().UnixNano()
	for i := 0; i < n; i++ {
		dt := make([]int64, 0, m)
		d = append(d, dt)

		go func(i int) {

			var reply proto.BenchmarkMessage

			startWg.Done()
			startWg.Wait()

			for j := 0; j < m; j++ {
				// 限流，这里不把限流的时间计算到等待耗时中
				tb.Wait(1)

				t := time.Now().UnixNano()
				ci := atomic.AddUint64(&clientIndex, 1)
				ci = ci % uint64(*pool)
				xclient := poolClients[int(ci)]

				err := xclient.Call(context.Background(), serviceMethod, args, &reply)
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

	// 等待测试完成
	wg.Wait()

	// 统计
	totalTInNano := time.Now().UnixNano() - totalT
	totalT = totalTInNano / 1000000
	log.Infof("took %d ms for %d requests", totalT, n*m)

	totalD := make([]int64, 0, n*m)
	for _, k := range d {
		totalD = append(totalD, k...)
	}
	totalD2 := make([]float64, 0, n*m)
	for _, k := range totalD {
		totalD2 = append(totalD2, float64(k))
	}

	mean, _ := stats.Mean(totalD2)
	median, _ := stats.Median(totalD2)
	max, _ := stats.Max(totalD2)
	min, _ := stats.Min(totalD2)
	p995, _ := stats.Percentile(totalD2, 99.5)

	log.Infof("sent     requests    : %d\n", n*m)
	log.Infof("received requests    : %d\n", atomic.LoadUint64(&trans))
	log.Infof("received requests_OK : %d\n", atomic.LoadUint64(&transOK))
	if totalT == 0 {
		log.Infof("throughput  (TPS)    : %d\n", int64(n*m)*1000*1000000/totalTInNano)
	} else {
		log.Infof("throughput  (TPS)    : %d\n", int64(n*m)*1000/totalT)
	}

	log.Infof("mean: %.f ns, median: %.f ns, max: %.f ns, min: %.f ns, p99.5: %.f ns\n", mean, median, max, min, p995)
	log.Infof("mean: %d ms, median: %d ms, max: %d ms, min: %d ms, p99.5: %d ms\n", int64(mean/1000000), int64(median/1000000), int64(max/1000000), int64(min/1000000), int64(p995/1000000))

}

// 准备请求数据
func prepareArgs() *proto.BenchmarkMessage {
	b := true
	var i int32 = 100000
	var s = "许多往事在眼前一幕一幕，变的那麼模糊"

	var args proto.BenchmarkMessage

	v := reflect.ValueOf(&args).Elem()
	num := v.NumField()
	for k := 0; k < num; k++ {
		field := v.Field(k)
		if field.Type().Kind() == reflect.Ptr {
			switch v.Field(k).Type().Elem().Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64:
				field.Set(reflect.ValueOf(&i))
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
