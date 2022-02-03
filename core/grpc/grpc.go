package graph_grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/jeconias/graph-service/graphpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	graph_database "github.com/jeconias/graph-service/core/database"
	"github.com/jeconias/graph-service/core/database/schema"
)

type GrpcServer struct {
	database *graph_database.Database
	Close    func()
}

func (v *GrpcServer) Init(db *graph_database.Database) {
	v.database = db
}

func (v *GrpcServer) Run(address string) {
	listen, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Net Listen Error: %v", err)
	}

	s := grpc.NewServer()
	graphpb.RegisterGraphServiceServer(s, v)

	reflection.Register(s)

	v.Close = func() {
		fmt.Println("Stopped GRPC Server")
		s.Stop()

		if err = v.database.Conn.Disconnect(context.TODO()); err != nil {
			log.Fatalf("failed to disconnect Database: %v", err)
		}

	}

	fmt.Println("Running GRPC Server on port 50051")

	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func (v *GrpcServer) Register(ctx context.Context, req *graphpb.GraphRequest) (*graphpb.GraphResponse, error) {

	resp, err := v.database.InsertVertice(schema.VerticeSchema{
		From: req.From,
		To:   req.To,
		Infos: []schema.VerticeInfoSchema{
			{
				Url:  req.Infos.Url,
				Date: req.Infos.Date,
			},
		},
	})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed on UpsertVertice: %s\n", err))
	}

	return &graphpb.GraphResponse{
		Message: "Success",
		Data: &graphpb.Node{
			Id:   resp.ID.Hex(),
			From: resp.From,
			To:   resp.To,
		},
	}, nil
}

func (v *GrpcServer) RegisterStream(stream graphpb.GraphService_RegisterStreamServer) error {

	readyRows := make([]schema.VerticeSchema, 0, 5)
	wait := make(chan string)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			go v.RegisterStreamUpsert(readyRows, wait)
			return nil
		}

		if err != nil {
			return errors.New(fmt.Sprintf("Error while reading client stream: %s", err))
		}

		if len(readyRows) >= 5 {
			data := readyRows
			readyRows = readyRows[:0]

			fmt.Println("Start write data...")
			go v.RegisterStreamUpsert(data, wait)

			fmt.Printf("Channel Status: %s\n", <-wait)
		}

		fmt.Println("Reading data...")
		readyRows = append(readyRows, schema.VerticeSchema{
			From: req.From,
			To:   req.To,
			Infos: []schema.VerticeInfoSchema{
				{
					Url:  req.Infos.Url,
					Date: req.Infos.Date,
				},
			},
		})

	}
}

func (v *GrpcServer) RegisterStreamUpsert(data []schema.VerticeSchema, wait chan string) error {

	_, dbErr := v.database.UpsertManyVertice(data)

	wait <- "Now you can read more data"

	if dbErr != nil {
		fmt.Printf("Error while writing client data: %s\n", dbErr)
		return dbErr
	}

	fmt.Printf("Success on writing client data: %d rows\n", len(data))

	return nil
}
