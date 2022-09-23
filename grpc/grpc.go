package grpc

import (
	"context"
	pb "int-service/_proto"
	"int-service/service"

	"github.com/sirupsen/logrus"
)

type GrpcServer struct {
	logger  *logrus.Logger
	service service.Servicer
	pb.UnimplementedClothingSvcServer
}

type GrpcServerProject struct {
	logger  *logrus.Logger
	service service.ProjectServicer
	pb.UnimplementedArticleSvcServer
	pb.UnimplementedCelebritySvcServer
	pb.UnimplementedEpisodeSvcServer
	pb.UnimplementedShowSvcServer
	pb.UnimplementedGenreSvcServer
	pb.UnimplementedSeasonSvcServer
	pb.UnimplementedJournalistSvcServer
}

func New(service service.Servicer, logger *logrus.Logger) *GrpcServer {
	return &GrpcServer{
		logger:  logger,
		service: service,
	}
}

func NewSvc(service service.ProjectServicer, logger *logrus.Logger) *GrpcServerProject {
	return &GrpcServerProject{
		logger:  logger,
		service: service,
	}
}

func (s *GrpcServer) CreateClothing(ctx context.Context, req *pb.CreateClothingRequest) (*pb.Clothing, error) {
	resp, err := s.service.CreateClothing(ctx, req.Type, int(req.Size), int(req.Price), req.Gender)
	if err != nil {
		return nil, err
	}

	return resp.ToGrpc().(*pb.Clothing), nil
}

func (s *GrpcServer) DeleteClothing(ctx context.Context, req *pb.DeleteClothingRequest) (*pb.EmptyResponse, error) {
	err := s.service.DeleteClothing(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (s *GrpcServer) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.ClothingListResponse, error) {
	resp, err := s.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	clothes := &pb.ClothingListResponse{}
	for _, c := range resp {
		clothes.Clothes = append(clothes.Clothes, c.ToGrpc().(*pb.Clothing))
	}
	return clothes, nil
}
