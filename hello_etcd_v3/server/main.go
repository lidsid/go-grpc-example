package main

import (
  "flag"
  "log"
  "net"
  "os"
  "os/signal"
  "syscall"

  pb "go-grpc-example/proto/hello"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "fmt"
  balancer "go-grpc-example/hello_etcd_v3"
)

const svcName = "project/test"

// Address gRPC服务地址
var addr = "127.0.0.1:50051"


// 定义helloService并实现约定的接口
type helloService struct{}

// HelloService Hello服务
var HelloService = helloService{}

// SayHello 实现Hello服务接口
func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
  resp := new(pb.HelloResponse)
  fmt.Println(in.Name)
  resp.Message = fmt.Sprintf("Hello %s.", in.Name)

  return resp, nil
}

func main() {
  flag.StringVar(&addr, "addr", addr, "addr to lis")
  flag.Parse()

  lis, err := net.Listen("tcp", addr)
  if err != nil {
    log.Fatalf("failed to listen: %s", err)
  }
  defer lis.Close()

  s := grpc.NewServer()
  defer s.GracefulStop()

  pb.RegisterHelloServer(s, &helloService{})

  go balancer.Register("127.0.0.1:2379", svcName, addr, 5)

  ch := make(chan os.Signal, 1)
  signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
  go func() {
    s := <-ch
    balancer.UnRegister(svcName, addr)

    if i, ok := s.(syscall.Signal); ok {
      os.Exit(int(i))
    } else {
      os.Exit(0)
    }

  }()

  if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %s", err)
  }
}