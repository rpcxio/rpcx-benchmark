module github.com/rpcxio/rpcx-benchmark/kitex

go 1.18

require (
	github.com/cloudwego/kitex v0.3.4
	github.com/gogo/protobuf v1.3.2
	github.com/rpcxio/rpcx-benchmark/proto v0.0.0-20220730153134-292b05d6ec48
	github.com/rpcxio/rpcx-benchmark/stat v0.0.0-20220730084343-905328fa1a4b
	github.com/smallnest/rpcx v1.7.8
	go.uber.org/ratelimit v0.2.0
	golang.org/x/net v0.0.0-20220728211354-c7608f3a8462
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/apache/thrift v0.16.0 // indirect
	github.com/bytedance/gopkg v0.0.0-20220728082022-984d38a4af78 // indirect
	github.com/chenzhuoyu/iasm v0.0.0-20220520152703-997ea6739ce9 // indirect
	github.com/choleraehyq/pid v0.0.14 // indirect
	github.com/cloudwego/frugal v0.1.1 // indirect
	github.com/cloudwego/netpoll v0.2.5 // indirect
	github.com/cloudwego/thriftgo v0.1.7 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.6.6 // indirect
	github.com/oleiade/lane v1.0.1 // indirect
	github.com/tidwall/gjson v1.14.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	golang.org/x/arch v0.0.0-20220722155209-00200b7164a7 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220728213248-dd149ef739b9 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0

replace github.com/rpcxio/rpcx-benchmark/proto => ../proto
