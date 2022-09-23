package grpc

import (
	"context"
	"errors"
	pb "int-service/_proto"
	"int-service/models"
	"time"
)

func (s *GrpcServerProject) CreateShow(ctx context.Context, req *pb.CreateShowRequest) (*pb.Show, error) {
	releaseDate := req.ReleaseDate.AsTime()
	if err := req.ReleaseDate.CheckValid(); err != nil {
		return nil, errors.New("error while converting releaseDate")
	}
	var endDate time.Time
	if req.EndDate != nil {
		endDate = req.EndDate.AsTime()
		if err := req.EndDate.CheckValid(); err != nil {
			return nil, errors.New("error while converting endDate")
		}
	}

	length := &models.ShowLength{
		Hours:   int(req.Length.Hours),
		Minutes: int(req.Length.Minutes),
	}

	resp, err := s.service.CreateShow(ctx, req.Title, req.Type, req.PostersPath, releaseDate, endDate, req.Rating, length, req.TrailerUrl, toShortGenresModel(req.Genres), toFilmCrewsModel(req.DirectedBy), toFilmCrewsModel(req.ProducedBy), toFilmCrewsModel(req.WrittenBy), toShortCelebsModel(req.Starring), req.Description, toShortSeasonsModel(req.Seasons))
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Show), nil
}

func (s *GrpcServerProject) GetShow(ctx context.Context, req *pb.GetByIDRequest) (*pb.Show, error) {
	resp, err := s.service.GetShow(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Show), nil
}

func (s *GrpcServerProject) UpdateShow(ctx context.Context, req *pb.Show) (*pb.Show, error) {
	releaseDate := req.ReleaseDate.AsTime()
	if err := req.ReleaseDate.CheckValid(); err != nil {
		return nil, errors.New("error while converting releaseDate")
	}
	var endDate time.Time
	if req.EndDate != nil {
		if err := req.EndDate.CheckValid(); err != nil {
			return nil, errors.New("error while converting endDate")
		}
	}

	length := &models.ShowLength{
		Hours:   int(req.Length.Hours),
		Minutes: int(req.Length.Minutes),
	}
	resp, err := s.service.UpdateShow(ctx, req.Id, req.Title, req.Type, req.PostersPath, releaseDate, endDate, req.Rating, length, req.TrailerUrl, toShortGenresModel(req.Genres), toFilmCrewsModel(req.DirectedBy), toFilmCrewsModel(req.ProducedBy), toFilmCrewsModel(req.WrittenBy), toShortCelebsModel(req.Starring), req.Description, toShortSeasonsModel(req.Seasons))
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Show), nil
}

func (s *GrpcServerProject) ListShows(ctx context.Context, req *pb.GetAllRequest) (*pb.ShowListResponse, error) {
	resp, err := s.service.ListShows(ctx)
	if err != nil {
		return nil, err
	}

	shows := &pb.ShowListResponse{}
	for _, show := range resp {
		shows.Shows = append(shows.Shows, show.ToGrpc().(*pb.Show))
	}
	return shows, nil
}

func (s *GrpcServerProject) UploadSeriesPosters(ctx context.Context, req *pb.UploadSeriesPostersRequest) (*pb.Show, error) {
	resp, err := s.service.UploadSeriesPosters(ctx, req.SeriesId, req.PostersPath)
	if err != nil {
		return nil, errors.New("error while updating series posters")
	}
	return resp.ToGrpc().(*pb.Show), nil
}

func (s *GrpcServerProject) DeleteSeriesPoster(ctx context.Context, req *pb.DeleteSeriesPosterRequest) (*pb.EmptyResponse, error) {
	err := s.service.DeleteSeriesPoster(ctx, req.SeriesId, req.Image)
	if err != nil {
		return nil, errors.New("error while updating deleted series poster")
	}
	return &pb.EmptyResponse{}, nil
}

func (s *GrpcServerProject) UploadMoviePosters(ctx context.Context, req *pb.UploadMoviePostersRequest) (*pb.Show, error) {
	resp, err := s.service.UploadMoviePosters(ctx, req.MovieId, req.PostersPath)
	if err != nil {
		return nil, errors.New("error while updating movie posters")
	}
	return resp.ToGrpc().(*pb.Show), nil
}

func (s *GrpcServerProject) DeleteMoviePoster(ctx context.Context, req *pb.DeleteMoviePosterRequest) (*pb.EmptyResponse, error) {
	err := s.service.DeleteMoviePoster(ctx, req.MovieId, req.Image)
	if err != nil {
		return nil, errors.New("error while updating deleted movie poster")
	}
	return &pb.EmptyResponse{}, nil
}