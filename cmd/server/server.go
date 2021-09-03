package main

import (
	"context"
	"tf-grpc-svc/internal/runner"
	"tf-grpc-svc/pkg/pb"
)

// Prove that Server implements pb.CalculatorServer by instantiating a Server
var _ pb.TerraformServer = (*Server)(nil)

type Server struct {
	pb.UnimplementedTerraformServer
	runner.TfRunner
}

func (s *Server) Validate(context.Context, *pb.NoInput) (*pb.ValidateResult, error) {
	return nil, nil
}

func (s *Server) Plan(in *pb.NoInput, stream pb.Terraform_PlanServer) error {
	outputSender := func(chunk string) {
		err := stream.Send(&pb.CommandOutput{
			Chunk: chunk,
		})

		if err != nil {
			return
		}
	}

	err := s.TfRunner.Plan(stream.Context(), "/home/mitchell/projects/tt/terraform", outputSender)

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Apply(in *pb.NoInput, stream pb.Terraform_ApplyServer) error {
	return nil
}

func NewServer() *Server {
	return &Server{
		TfRunner: *runner.NewRunner("/usr/bin/terraform"),
	}
}
