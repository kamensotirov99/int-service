package grpc

import (
	"context"
	pb "int-service/_proto"
	"int-service/models"
)

func (s *GrpcServerProject) CreateGenre(ctx context.Context, req *pb.CreateGenreRequest) (*pb.Genre, error) {
	resp, err := s.service.CreateGenre(ctx, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Genre), nil
}

func (s *GrpcServerProject) GetGenre(ctx context.Context, req *pb.GetByIDRequest) (*pb.Genre, error) {
	resp, err := s.service.GetGenre(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Genre), nil
}

func (s *GrpcServerProject) GetGenreByName(ctx context.Context, req *pb.GetByNameRequest) (*pb.Genre, error) {
	resp, err := s.service.GetGenreByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Genre), nil
}

func (s *GrpcServerProject) UpdateGenre(ctx context.Context, req *pb.Genre) (*pb.Genre, error) {
	resp, err := s.service.UpdateGenre(ctx, req.Id, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Genre), nil
}

func (s *GrpcServerProject) ListGenres(ctx context.Context, req *pb.GetAllRequest) (*pb.GenreListResponse, error) {
	resp, err := s.service.ListGenres(ctx)
	if err != nil {
		return nil, err
	}

	genres := &pb.GenreListResponse{}
	for _, genre := range resp {
		genres.Genres = append(genres.Genres, genre.ToGrpc().(*pb.Genre))
	}
	return genres, nil
}

func toShortGenresModel(genresPb *pb.ShortGenres) models.ShortGenres {
	genres := models.ShortGenres{}
	for _, genre := range genresPb.Genres {
		genres = append(genres, &models.ShortGenre{
			ID:          genre.Id,
			Name:        genre.Name,
		})
	}
	return genres
}
