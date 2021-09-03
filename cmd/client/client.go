package main

import (
	"google.golang.org/grpc"
	"tf-grpc-svc/pkg/pb"
)

type Client struct {
	client pb.TerraformClient
}

func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		client: pb.NewTerraformClient(conn),
	}
}
