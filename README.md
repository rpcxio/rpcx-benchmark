# 流行的RPC框架测试

先前的[rpcx-benchmark](https://github.com/rpcx-ecosystem/rpcx-benchmark)提供了一个很好的各种RPC框架性能测试结果。它采用启动一定数量的线程/纤程，每个线程独立使用一个client进行服务调用，用来获取服务调用的吞吐率和耗时。这种测试模式模拟了现实生活中很多的业务客户端并发访问RPC服务的时候，RPC服务所能提供的能力(吞吐率和延迟，这是一对相爱相杀的性能指标)。

但是，测试最难的是无法实现一个“银弹”，来模拟现实中所有的场景，所以新建了这个测试工程，期望能测试更多的场景：

- 能够在相同的吞吐率情况下比较耗时。 但是需要注意的一点是，有些RPC框架并不能达到很高的吞吐率，在耗时还比较合理的情况下，我们不能迁就低吞吐率的框架，所以会在测试结果中标记出来。
- 采用共享的client。创建一定数量的client,在线程/纤程中共享。用来测试单个业务能采用对象池的方式。如果你想了解不共享的情况，请参考先前的测试项目[rpcx-benchmark](https://github.com/rpcx-ecosystem/rpcx-benchmark)。
- 所有的框架都是在“公平”的情况下测试。测试数据都是一致的，采用protobuf进行测试。虽然有比Protobuf性能更好的序列化框架，但是因为不具有通用性所以不考虑。
- 测试会进行预热。
- 避免[coordinated omission](http://highscalability.com/blog/2015/10/5/your-load-generator-is-probably-lying-to-you-take-the-red-pi.html):测试统计的是`等待时间`+`服务时间`,而不是服务端服务时间
- 统计既包含平均值，也包含中位值、最大值、最小值和 P99值。

每个rpc框架包含一个单独的子文件夹，里面包含一个server文件夹和client文件夹。server文件用来启动一个server，client文件夹用来测试并输出测试结果。

你可以参考已存在的框架测试代码添加你相测试的RPC框架，客户端和服务器端必须要包含相同的序列化和反序列化逻辑，以及相同的处理逻辑，不同通过hack的方式绕过序列化反序列化违反公平性。

`proto`包含测试用的`proto`文件和生成的Go代码文件，这是公用的。如果需要，你也可以生成自己的Go代码，但是需要基于同样的`proto`文件。如果你修改成`proto3`的代码，我认为也是合理的和可以接受的。

**注意**

- 不要使用并发数1作为最终测试结果，而是仅仅用来调试程序。 因为实际项目中很少单个客户端访问服务端的场景，如果有的话，那么你需要考虑是否有必要拆分成微服务。
- 并发数多测试几个场景，比如并发客户端为500、1000、2000、5000、10000的场景。rpc client对象池要小于并发客户端数，用来测试客户端内部可以共享rpc client的场景。
- 耗时太大的情况下比较吞吐率是没有意义的。没有实际业务会在服务几乎不可用的情况下还谈论吞吐率。
- 请求错误数太多的情况下指标是没有意义的。例如所有的请求都返回错误，耗时可能很低，但这是无意义的。


## RPC 框架

测试的RPC框架包括:

- [rpcx](https://github.com/rpcxio/rpcx-benchmark/tree/master/rpcx)
- [grpc/go](https://github.com/rpcxio/rpcx-benchmark/tree/master/grpc)
- [Go标准库的rpc](https://github.com/rpcxio/rpcx-benchmark/tree/master/go_stdrpc)
- [tarsgo](): todo
- [thrift](): todo
- [dubbo-go](): todo
- [hprose](): todo
- [go-micro](): toto

欢迎补充`todo`的代码，欢迎提交其它rpc框架的测试代码。