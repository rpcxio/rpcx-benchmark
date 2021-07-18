package benchmark

import (
	"time"

	"github.com/montanaflynn/stats"
	"github.com/smallnest/rpcx/log"
)

// Stats 统计结果.
func Stats(startTime int64, totalRequests int, tookTimes [][]int64, trans, transOK uint64) {
	// 测试总耗时
	totalTInNano := time.Now().UnixNano() - startTime
	totalT := totalTInNano / 1000000
	log.Infof("took %d ms for %d requests", totalT, totalRequests)

	// 汇总每个请求的耗时
	totalD := make([]int64, 0, totalRequests)
	for _, k := range tookTimes {
		totalD = append(totalD, k...)
	}
	// 将int64数组转换成float64数组，以便分析
	totalD2 := make([]float64, 0, totalRequests)
	for _, k := range totalD {
		totalD2 = append(totalD2, float64(k))
	}

	// 计算各个指标
	mean, _ := stats.Mean(totalD2)
	median, _ := stats.Median(totalD2)
	max, _ := stats.Max(totalD2)
	min, _ := stats.Min(totalD2)
	p999, _ := stats.Percentile(totalD2, 99.9)

	// 输出结果
	log.Infof("sent     requests    : %d\n", totalRequests)
	log.Infof("received requests    : %d\n", trans)
	log.Infof("received requests_OK : %d\n", transOK)
	if totalT == 0 {
		log.Infof("throughput  (TPS)    : %d\n", int64(totalRequests)*1000*1000000/totalTInNano)
	} else {
		log.Infof("throughput  (TPS)    : %d\n\n", int64(totalRequests)*1000/totalT)
	}

	log.Infof("mean: %.f ns, median: %.f ns, max: %.f ns, min: %.f ns, p99.9: %.f ns\n", mean, median, max, min, p999)
	log.Infof("mean: %d ms, median: %d ms, max: %d ms, min: %d ms, p99.9: %d ms\n", int64(mean/1000000), int64(median/1000000), int64(max/1000000), int64(min/1000000), int64(p999/1000000))
}
