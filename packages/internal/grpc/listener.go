package grpc

import (
	"log"
	"net"
	"sync"

	"example.com/go-basic/proto"
	"google.golang.org/grpc"
)

func RunGrpcListener(grpcSrv *GrpcServer) {
	lis, err := net.Listen("tcp", "localhost:50001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterEntityServiceServer(grpcServer, grpcSrv)

	log.Println("grpc server is running at :50001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func RunGrpcListenerInParallel(grpcSrv *GrpcServer) {
	var wg sync.WaitGroup

	wg.Add(1)
	defer wg.Done()
	RunGrpcListener(grpcSrv)
	wg.Wait()
}
