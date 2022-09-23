package grpc

import (
	"context"
	pb "int-service/_proto"
)

func (s *GrpcServerProject) CreateJournalist(ctx context.Context, req *pb.CreateJournalistRequest) (*pb.Journalist, error) {
	resp, err := s.service.CreateJournalist(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Journalist), nil
}

func (s *GrpcServerProject) GetJournalist(ctx context.Context, req *pb.GetByIDRequest) (*pb.Journalist, error) {
	resp, err := s.service.GetJournalist(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Journalist), nil
}

func (s *GrpcServerProject) GetJournalistByName(ctx context.Context, req *pb.GetByNameRequest) (*pb.Journalist, error) {
	resp, err := s.service.GetJournalistByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Journalist), nil
}

func (s *GrpcServerProject) UpdateJournalist(ctx context.Context, req *pb.Journalist) (*pb.Journalist, error) {
	resp, err := s.service.UpdateJournalist(ctx, req.Id, req.Name)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Journalist), nil
}

func (s *GrpcServerProject) ListJournalists(ctx context.Context, req *pb.GetAllRequest) (*pb.JournalistListResponse, error) {
	resp, err := s.service.ListJournalists(ctx)
	if err != nil {
		return nil, err
	}

	journalists := &pb.JournalistListResponse{}
	for _, j := range resp {
		journalists.Journalists = append(journalists.Journalists, j.ToGrpc().(*pb.Journalist))
	}
	return journalists, nil
}
