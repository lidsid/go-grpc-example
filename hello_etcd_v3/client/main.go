package main

import (
  "fmt"
  "time"

  pb "go-grpc-example/proto/hello"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "google.golang.org/grpc/resolver"
  balancer "go-grpc-example/hello_etcd_v3"
)

func main() {
  r := balancer.NewResolver("localhost:2379")
  resolver.Register(r)

  conn, err := grpc.Dial(r.Scheme()+"://author/project/test", grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), grpc.WithInsecure())
  if err != nil {
    panic(err)
  }

  client := pb.NewHelloClient(conn)

  for {
    resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "hello"}, grpc.WaitForReady(true))
    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Println(resp)
    }

    <-time.After(time.Second)
  }
}