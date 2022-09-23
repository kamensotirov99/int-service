package grpc

import (
	"context"
	"errors"
	pb "int-service/_proto"
	"int-service/models"
)

func (s *GrpcServerProject) CreateSeason(ctx context.Context, req *pb.CreateSeasonRequest) (*pb.Season, error) {
	releaseDate := req.ReleaseDate.AsTime()
	if err := req.ReleaseDate.CheckValid(); err != nil {
		return nil, errors.New("error while converting releaseDate")
	}
	resp, err := s.service.CreateSeason(ctx, req.ShowId, req.Title, req.TrailerUrl, []string{}, releaseDate, req.Rating, req.Resume, toFilmCrewsModel(req.DirectedBy), toFilmCrewsModel(req.ProducedBy), toFilmCrewsModel(req.WrittenBy), toShortEpisodesModel(req.Episodes))
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Season), nil
}

func (s *GrpcServerProject) GetSeason(ctx context.Context, req *pb.GetByIDRequest) (*pb.Season, error) {
	resp, err := s.service.GetSeason(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Season), nil
}

func (s *GrpcServerProject) UpdateSeason(ctx context.Context, req *pb.Season) (*pb.Season, error) {
	releaseDate := req.ReleaseDate.AsTime()
	if err := req.ReleaseDate.CheckValid(); err != nil {
		return nil, errors.New("error while converting releaseDate")
	}
	resp, err := s.service.UpdateSeason(ctx, req.Id, req.ShowId, req.Title, req.TrailerUrl, req.PostersPath, releaseDate, req.Rating, req.Resume, toFilmCrewsModel(req.DirectedBy), toFilmCrewsModel(req.ProducedBy), toFilmCrewsModel(req.WrittenBy), toShortEpisodesModel(req.Episodes))
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Season), nil
}

func (s *GrpcServerProject) UploadSeasonPosters(ctx context.Context, req *pb.UploadSeasonPostersRequest) (*pb.Season, error) {
	resp, err := s.service.UploadSeasonPosters(ctx, req.SeasonId, req.PostersPath)
	if err != nil {
		return nil, errors.New("error while uploading season posters")
	}
	return resp.ToGrpc().(*pb.Season), nil
}

func (s *GrpcServerProject) DeleteSeasonPoster(ctx context.Context, req *pb.DeleteSeasonPosterRequest) (*pb.EmptyResponse, error) {
	err := s.service.DeleteSeasonPoster(ctx, req.SeriesId, req.SeasonId, req.Image)
	if err != nil {
		return nil, errors.New("error while deleting season poster")
	}
	return &pb.EmptyResponse{}, nil
}

func (s *GrpcServerProject) ListShowSeasons(ctx context.Context, req *pb.GetByIDRequest) (*pb.ListSeasonResponse, error) {
	resp, err := s.service.ListShowSeasons(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	showSeasons := &pb.ListSeasonResponse{}
	for _, season := range resp {
		showSeasons.Seasons = append(showSeasons.Seasons, season.ToGrpc().(*pb.Season))
	}
	return showSeasons, nil
}

func (s *GrpcServerProject) ListSeasonsCollection(ctx context.Context, req *pb.GetAllRequest) (*pb.ListSeasonResponse, error) {
	resp, err := s.service.ListSeasonsCollection(ctx)
	if err != nil {
		return nil, err
	}
	collection := &pb.ListSeasonResponse{}
	for _, season := range resp {
		collection.Seasons = append(collection.Seasons, season.ToGrpc().(*pb.Season))
	}
	return collection, nil
}

func toShortSeasonsModel(seasonsPb *pb.ShortSeasons) models.ShortSeasons {
	shortSeasons := models.ShortSeasons{}
	for _, season := range seasonsPb.Seasons {
		shortSeasons = append(shortSeasons, &models.ShortSeason{
			ID:          season.Id,
			Title:       season.Title,
			PostersPath: season.PostersPath,
			Rating:      season.Rating,
		})
	}
	return shortSeasons
}
