module github.com/rpcxio/rpcx-benchmark/kitex

go 1.16

require (
	github.com/cloudwego/kitex v0.0.2-0.20210722063233-67262666003d
	github.com/fatih/color v1.12.0 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/rpcxio/rpcx-benchmark/stat v0.0.0-20210725144741-82e06e86a2a9
	github.com/smallnest/rpcx v1.6.4
	go.uber.org/ratelimit v0.2.0
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/protobuf v1.26.0
)

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0