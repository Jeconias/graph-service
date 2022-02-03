package main

import (
	"log"
	"os"
	"os/signal"

	graph_database "github.com/jeconias/graph-service/core/database"
	graph_grpc "github.com/jeconias/graph-service/core/grpc"
	graph_http "github.com/jeconias/graph-service/core/http"
)

func main() {

	graphDatabase, err := graph_database.InitDB()
	if err != nil {
		log.Fatalf("Database not started: %s\n", err)
		return
	}

	graphHttp := graph_http.GraphHttp{}
	grpcServer := graph_grpc.GrpcServer{}

	graphHttp.Init(graphDatabase)
	grpcServer.Init(graphDatabase)

	go func() {
		grpcServer.Run(":50051")
		graphHttp.Run(":3001")
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	grpcServer.Close()
}
