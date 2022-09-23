package service

import (
	"context"
	"int-service/dto"
	"int-service/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type EpisodeServicer interface {
	CreateEpisode(ctx context.Context, seasonID string, title string, postersPath []string, trailerURL string, length *models.ShowLength, rating float64, resume string, writtenBy models.FilmCrews, producedBy models.FilmCrews, directedBy models.FilmCrews, starring models.ShortCelebrities) (models.ResponseModeler, error)
	GetEpisode(ctx context.Context, ID string) (models.ResponseModeler, error)
	UpdateEpisode(ctx context.Context, ID string, seasonID string, title string, postersPath []string, trailerURL string, length *models.ShowLength, rating float64, resume string, writtenBy models.FilmCrews, producedBy models.FilmCrews, directedBy models.FilmCrews, starring models.ShortCelebrities) (models.ResponseModeler, error)
	UploadEpisodePosters(ctx context.Context, episodeID string, postersPath []string) (models.ResponseModeler, error)
	DeleteEpisodePoster(ctx context.Context, seriesID string, seasonID string, episodeID string, image string) error
	ListSeasonEpisodes(ctx context.Context, seasonID string) ([]models.ResponseModeler, error)
	ListCollectionEpisodes(ctx context.Context) ([]models.ResponseModeler, error)
}

func (s *projectService) CreateEpisode(ctx context.Context, seasonID string, title string, postersPath []string, trailerURL string, length *models.ShowLength, rating float64, resume string, writtenBy models.FilmCrews, producedBy models.FilmCrews, directedBy models.FilmCrews, starring models.ShortCelebrities) (models.ResponseModeler, error) {
	err := s.validateEpisodeUniqueness(ctx, seasonID, title)
	if err != nil {
		s.logger.Error("Error while creating episode")
		return nil, errors.Wrap(err, "Error while creating episode")
	}
	episode := toEpisodeDTO(uuid.New().String(), seasonID, title, postersPath, trailerURL, length, rating, resume, writtenBy, producedBy, directedBy, starring)
	resp, err := s.repository.CreateEpisode(ctx, episode)
	if err != nil {
		s.logger.Error("Error while creating episode")
		return nil, errors.Wrap(err, "Error while creating episode")
	}
	_, err = s.repository.AddShortEpisode(ctx, seasonID, &dto.ShortEpisodeDTO{
		ID:          episode.ID,
		Title:       episode.Title,
		PostersPath: episode.PostersPath,
		Rating:      episode.Rating,
		Resume:      episode.Resume,
	})
	if err != nil {
		s.logger.Error("Error while adding short episode in season")
		return nil, errors.Wrap(err, "Error while adding short episode in season")
	}
	return resp.ToModel(), nil
}

func (s *projectService) GetEpisode(ctx context.Context, ID string) (models.ResponseModeler, error) {
	resp, err := s.repository.GetEpisode(ctx, ID)
	if err != nil {
		s.logger.Error("Error while getting episode by id")
		return nil, errors.Wrap(err, "Error while getting episode by id")
	}
	return resp.ToModel(), nil
}

func (s *projectService) UpdateEpisode(ctx context.Context, ID string, seasonID string, title string, postersPath []string, trailerURL string, length *models.ShowLength, rating float64, resume string, writtenBy models.FilmCrews, producedBy models.FilmCrews, directedBy models.FilmCrews, starring models.ShortCelebrities) (models.ResponseModeler, error) {
	updatedEpisode := toEpisodeDTO(ID, seasonID, title, postersPath, trailerURL, length, rating, resume, writtenBy, producedBy, directedBy, starring)
	resp, err := s.repository.UpdateEpisode(ctx, updatedEpisode)
	if err != nil {
		s.logger.Error("Error while updating episode")
		return nil, errors.New("Error while updating episode")
	}
	_, err = s.repository.UpdateShortEpisode(ctx, &dto.ShortEpisodeDTO{
		ID:          resp.ID,
		Title:       resp.Title,
		PostersPath: resp.PostersPath,
		Rating:      resp.Rating,
		Resume:      resp.Resume,
	})
	if err != nil {
		s.logger.Error("Error while updating short episode")
		return nil, errors.New("Error while updating short episode")
	}
	return resp.ToModel(), nil
}

func (s *projectService) UploadEpisodePosters(ctx context.Context, episodeID string, postersPath []string) (models.ResponseModeler, error) {
	resp, err := s.repository.UploadEpisodePosters(ctx, episodeID, postersPath)
	if err != nil {
		s.logger.Error("Error while uploading episode posters")
		return nil, errors.Wrap(err, "Error while uploading episode posters")
	}

	episode, err := s.repository.GetEpisode(ctx, episodeID)
	if err != nil {
		s.logger.Error("Error while getting episode by id")
		return nil, errors.Wrap(err, "Error while getting episode by id")
	}
	shortEpisode := &dto.ShortEpisodeDTO{
		ID:          episodeID,
		Title:       episode.Title,
		Rating:      episode.Rating,
		PostersPath: episode.PostersPath,
		Resume:      episode.Resume,
	}
	_, err = s.repository.UpdateShortEpisode(ctx, shortEpisode)
	if err != nil {
		s.logger.Error("Error while updating short episode posters")
		return nil, errors.Wrap(err, "Error while updating short episode posters")
	}
	return resp.ToModel(), nil
}

func (s *projectService) DeleteEpisodePoster(ctx context.Context, seriesID string, seasonID string, episodeID string, image string) error {
	err := s.repository.DeleteEpisodePoster(ctx, seriesID, seasonID, episodeID, image)
	if err != nil {
		s.logger.Error("Error while deleting episode poster in database")
		return errors.Wrap(err, "Error while deleting episode poster in database")
	}
	return nil
}

func (s *projectService) ListSeasonEpisodes(ctx context.Context, seasonID string) ([]models.ResponseModeler, error) {
	resp, err := s.repository.ListSeasonEpisodes(ctx, seasonID)
	if err != nil {
		s.logger.Error("Error while listing all season episodes")
		return nil, errors.New("Error while listing all season episodes")
	}
	episodes := []models.ResponseModeler{}
	for _, episode := range resp {
		episodes = append(episodes, episode.ToModel())
	}
	return episodes, nil
}

func (s *projectService) ListCollectionEpisodes(ctx context.Context) ([]models.ResponseModeler, error) {
	resp, err := s.repository.ListCollectionEpisodes(ctx)
	if err != nil {
		s.logger.Error("Error while listing all episodes")
		return nil, errors.New("Error while listing all episodes")
	}
	collection := []models.ResponseModeler{}
	for _, episode := range resp {
		collection = append(collection, episode.ToModel())
	}
	return collection, nil
}

func (s *projectService) validateEpisodeUniqueness(ctx context.Context, seasonID string, title string) error {
	episodes, err := s.repository.ListSeasonEpisodes(ctx, seasonID)
	if err != nil {
		return errors.Wrap(err, "Error while getting episodes from the Mongo database")
	}
	for _, episode := range episodes {
		if episode.Title == title {
			return errors.New("There is already an episode with that name in the Mongo database with id:" + episode.ID)
		}
	}
	return nil
}

func toEpisodeDTO(ID string, seasonID string, title string, postersPath []string, trailerURL string, lengthModel *models.ShowLength, rating float64, resume string, writtenBy models.FilmCrews, producedBy models.FilmCrews, directedBy models.FilmCrews, starring models.ShortCelebrities) *dto.EpisodeDTO {
	length := dto.ShowLengthDTO{
		Hours:   lengthModel.Hours,
		Minutes: lengthModel.Minutes,
	}
	episode := &dto.EpisodeDTO{
		ID:          ID,
		SeasonID:    seasonID,
		Title:       title,
		PostersPath: postersPath,
		TrailerURL:  trailerURL,
		Length:      length,
		Rating:      rating,
		Resume:      resume,
		WrittenBy:   toFilmCrewsDTO(writtenBy),
		ProducedBy:  toFilmCrewsDTO(producedBy),
		DirectedBy:  toFilmCrewsDTO(directedBy),
		Starring:    toShortCelebDTO(starring),
	}
	return episode
}

func toShortEpisodesDTO(episodes models.ShortEpisodes) dto.ShortEpisodesDTO {
	shortEpisodes := dto.ShortEpisodesDTO{}
	for _, episode := range episodes {
		shortEpisodes = append(shortEpisodes, &dto.ShortEpisodeDTO{
			ID:          episode.ID,
			Title:       episode.Title,
			PostersPath: episode.PostersPath,
			Rating:      episode.Rating,
			Resume:      episode.Resume,
		})
	}
	return shortEpisodes
}
