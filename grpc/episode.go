package grpc

import (
	"context"
	"errors"
	pb "int-service/_proto"
	"int-service/models"
)

func (s *GrpcServerProject) CreateEpisode(ctx context.Context, req *pb.CreateEpisodeRequest) (*pb.Episode, error) {
	length := &models.ShowLength{
		Hours:   int(req.ShowLength.Hours),
		Minutes: int(req.ShowLength.Minutes),
	}

	resp, err := s.service.CreateEpisode(ctx, req.SeasonId, req.Title, req.PostersPath, req.TrailerUrl, length, req.Rating, req.Resume, toFilmCrewsModel(req.WrittenBy), toFilmCrewsModel(req.ProducedBy), toFilmCrewsModel(req.DirectedBy), toShortCelebsModel(req.Starring))
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Episode), nil
}

func (s *GrpcServerProject) GetEpisode(ctx context.Context, req *pb.GetByIDRequest) (*pb.Episode, error) {
	resp, err := s.service.GetEpisode(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Episode), nil
}

func (s *GrpcServerProject) UpdateEpisode(ctx context.Context, req *pb.Episode) (*pb.Episode, error) {
	lenght := &models.ShowLength{
		Hours:   int(req.ShowLength.Hours),
		Minutes: int(req.ShowLength.Minutes),
	}
	resp, err := s.service.UpdateEpisode(ctx, req.Id, req.SeasonId, req.Title, req.PostersPath, req.TrailerUrl, lenght, req.Rating, req.Resume, toFilmCrewsModel(req.WrittenBy), toFilmCrewsModel(req.ProducedBy), toFilmCrewsModel(req.DirectedBy), toShortCelebsModel(req.Starring))
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Episode), nil
}

func (s *GrpcServerProject) UploadEpisodePosters(ctx context.Context, req *pb.UploadEpisodePostersRequest) (*pb.Episode, error) {
	resp, err := s.service.UploadEpisodePosters(ctx, req.EpisodeId, req.PostersPath)
	if err != nil {
		return nil, errors.New("error while updating episode posters")
	}
	return resp.ToGrpc().(*pb.Episode), nil
}

func (s *GrpcServerProject) DeleteEpisodePoster(ctx context.Context, req *pb.DeleteEpisodePosterRequest) (*pb.EmptyResponse, error) {
	err := s.service.DeleteEpisodePoster(ctx, req.SeriesId, req.SeasonId, req.EpisodeId, req.Image)
	if err != nil {
		return nil, errors.New("error while deleting episode poster")
	}
	return &pb.EmptyResponse{}, nil
}

func (s *GrpcServerProject) ListSeasonEpisodes(ctx context.Context, req *pb.GetByIDRequest) (*pb.ListEpisodeResponse, error) {
	resp, err := s.service.ListSeasonEpisodes(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	episodes := &pb.ListEpisodeResponse{}
	for _, episode := range resp {
		episodes.Episodes = append(episodes.Episodes, episode.ToGrpc().(*pb.Episode))
	}
	return episodes, nil
}

func (s *GrpcServerProject) ListCollectionEpisodes(ctx context.Context, req *pb.GetAllRequest) (*pb.ListEpisodeResponse, error) {
	resp, err := s.service.ListCollectionEpisodes(ctx)
	if err != nil {
		return nil, err
	}

	collection := &pb.ListEpisodeResponse{}
	for _, episodes := range resp {
		collection.Episodes = append(collection.Episodes, episodes.ToGrpc().(*pb.Episode))
	}
	return collection, nil
}

func toShortEpisodesModel(episodesPb *pb.ShortEpisodeList) models.ShortEpisodes {
	episodes := models.ShortEpisodes{}
	for _, episode := range episodesPb.ShortEpisodes {
		episodes = append(episodes, &models.ShortEpisode{
			ID:          episode.Id,
			Title:       episode.Title,
			PostersPath: episode.PostersPath,
			Rating:      episode.Rating,
			Resume:      episode.Resume,
		})
	}
	return episodes
}
