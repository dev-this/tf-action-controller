package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"tf-grpc-svc/pkg/pb"
)

var (
	host = flag.String("host", "localhost:8080", "gRPC server hostname")
)

type outputRecv interface {
	Recv() (*pb.CommandOutput, error)
}

func main() {
	flag.Parse()

	if *host == "" {
		log.Fatal("[main] unable to start client without gRPC endpoint to server")
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	var fireRequest func(ctx context.Context, client pb.TerraformClient) (outputRecv, error)

	switch cmd := flag.Arg(0); cmd {
	case "apply":
		fireRequest = func(ctx context.Context, client pb.TerraformClient) (outputRecv, error) {
			return client.Apply(ctx, &pb.NoInput{})
		}
	case "plan":
		fireRequest = func(ctx context.Context, client pb.TerraformClient) (outputRecv, error) {
			return client.Plan(ctx, &pb.NoInput{})
		}
	default:
		_, _ = fmt.Fprintf(os.Stderr, "unknown command '%s'\n", cmd)
		_, _ = fmt.Fprintln(os.Stderr, "commands available: plan, apply")
		os.Exit(1)
	}

	fmt.Println(flag.Arg(0))

	log.Printf("Connecting to gRPC Service [%s]", *host)
	conn, err := grpc.Dial(*host, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewTerraformClient(conn)
	ctx := context.Background()

	stream, err := fireRequest(ctx, client)
	if err != nil {
		// Panic so conn.Close() is invoked
		log.Panic(err)
	}

	err = handleStream(stream)

	if err != nil {
		log.Panicln(err)
	}
}

func handleStream(stream outputRecv) error {
	for {
		data, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		log.Println(data.Chunk)
	}

	return nil
}
