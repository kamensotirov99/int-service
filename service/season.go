package service

import (
	"context"
	"int-service/dto"
	"int-service/models"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type SeasonServicer interface {
	CreateSeason(ctx context.Context, showID string, title string, trailerURL string, postersPath []string, releaseDate time.Time, rating float64, resume string, directedBy models.FilmCrews, producedBy models.FilmCrews, writtenBy models.FilmCrews, episodes models.ShortEpisodes) (models.ResponseModeler, error)
	GetSeason(ctx context.Context, ID string) (models.ResponseModeler, error)
	UpdateSeason(ctx context.Context, ID string, showID string, title string, trailerURL string, postersPath []string, releaseDate time.Time, rating float64, resume string, directedBy models.FilmCrews, producedBy models.FilmCrews, writtenBy models.FilmCrews, episodes models.ShortEpisodes) (models.ResponseModeler, error)
	UploadSeasonPosters(ctx context.Context, seasonID string, postersPath []string) (models.ResponseModeler, error)
	DeleteSeasonPoster(ctx context.Context, seriesID string, seasonID string, image string) error
	ListShowSeasons(ctx context.Context, ID string) ([]models.ResponseModeler, error)
	ListSeasonsCollection(ctx context.Context) ([]models.ResponseModeler, error)
}

func (s *projectService) CreateSeason(ctx context.Context, showID string, title string, trailerURL string, postersPath []string, releaseDate time.Time, rating float64, resume string, directedBy models.FilmCrews, producedBy models.FilmCrews, writtenBy models.FilmCrews, episodes models.ShortEpisodes) (models.ResponseModeler, error) {
	err := s.validateSeasonUniqueness(ctx, showID, title)
	if err != nil {
		s.logger.Error("Error while creating season")
		return nil, errors.Wrap(err, "Error while creating season")
	}
	season := toSeasonDTO(uuid.New().String(), showID, title, trailerURL, postersPath, releaseDate, rating, resume, directedBy, producedBy, writtenBy, episodes)
	resp, err := s.repository.CreateSeason(ctx, season)
	if err != nil {
		s.logger.Error("Error while creating season")
		return nil, errors.Wrap(err, "Error while creating season")
	}
	_, err = s.repository.AddShortSeason(ctx, showID, &dto.ShortSeasonDTO{
		ID:          season.ID,
		Title:       season.Title,
		PostersPath: season.PostersPath,
		Rating:      season.Rating,
	})
	if err != nil {
		s.logger.Error("Error while adding short season in show")
		return nil, errors.Wrap(err, "Error while adding short season in show")
	}
	return resp.ToModel(), nil
}

func (s *projectService) GetSeason(ctx context.Context, ID string) (models.ResponseModeler, error) {
	resp, err := s.repository.GetSeason(ctx, ID)
	if err != nil {
		s.logger.Error("Error while getting season by id")
		return nil, errors.Wrap(err, "Error while getting season by id")
	}
	return resp.ToModel(), nil
}

func (s *projectService) UpdateSeason(ctx context.Context, ID string, showID string, title string, trailerURL string, postersPath []string, releaseDate time.Time, rating float64, resume string, directedBy models.FilmCrews, producedBy models.FilmCrews, writtenBy models.FilmCrews, episodes models.ShortEpisodes) (models.ResponseModeler, error) {
	updatedSeason := toSeasonDTO(ID, showID, title, trailerURL, postersPath, releaseDate, rating, resume, directedBy, producedBy, writtenBy, episodes)
	resp, err := s.repository.UpdateSeason(ctx, updatedSeason)
	if err != nil {
		s.logger.Error("Error while updating season")
		return nil, errors.New("Error while updating season")
	}
	_, err = s.repository.UpdateShortSeason(ctx, &dto.ShortSeasonDTO{
		ID:          resp.ID,
		Title:       resp.Title,
		PostersPath: resp.PostersPath,
		Rating:      resp.Rating,
	})

	if err != nil {
		s.logger.Error("Error while updating short season")
		return nil, errors.New("Error while updating short season")
	}

	return resp.ToModel(), nil
}

func (s *projectService) UploadSeasonPosters(ctx context.Context, seasonID string, postersPath []string) (models.ResponseModeler, error) {
	resp, err := s.repository.UploadSeasonPosters(ctx, seasonID, postersPath)
	if err != nil {
		s.logger.Error("Error while uploading season posters")
		return nil, errors.Wrap(err, "Error while uploading season posters")
	}

	season, err := s.repository.GetSeason(ctx, seasonID)
	if err != nil {
		s.logger.Error("Error while getting season by id")
		return nil, errors.Wrap(err, "Error while getting season by id")
	}
	shortSeason := &dto.ShortSeasonDTO{
		ID:          seasonID,
		Title:       season.Title,
		Rating:      season.Rating,
		PostersPath: season.PostersPath,
	}
	if _, err := s.repository.UpdateShortSeason(ctx, shortSeason); err != nil {
		s.logger.Error("Error while updating short season in show")
		return nil, errors.New("Error while updating short season in show")
	}

	return resp.ToModel(), nil
}

func (s *projectService) DeleteSeasonPoster(ctx context.Context, seriesID string, seasonID string, image string) error {
	err := s.repository.DeleteSeasonPoster(ctx, seriesID, seasonID, image)
	if err != nil {
		s.logger.Error("Error while deleting season poster in database")
		return errors.Wrap(err, "Error while deleting season poster in database")
	}

	if err := s.repository.DeleteShortSeasonPostersInShow(ctx, seriesID, seasonID, image); err != nil {
		s.logger.Error("Error while deleting short season poster in show")
		return errors.New("Error while deleting short season poster in show")
	}
	return nil
}

func (s *projectService) ListShowSeasons(ctx context.Context, ID string) ([]models.ResponseModeler, error) {
	resp, err := s.repository.ListShowSeasons(ctx, ID)
	if err != nil {
		s.logger.Error("Error while listing all show seasons")
		return nil, errors.New("Error while listing all show seasons")
	}
	seasons := []models.ResponseModeler{}
	for _, season := range resp {
		seasons = append(seasons, season.ToModel())
	}
	return seasons, nil
}

func (s *projectService) ListSeasonsCollection(ctx context.Context) ([]models.ResponseModeler, error) {
	resp, err := s.repository.ListSeasonsCollection(ctx)
	if err != nil {
		s.logger.Error("Error while listing all seasons")
		return nil, errors.New("Error while listing all seasons")
	}
	seasons := []models.ResponseModeler{}
	for _, season := range resp {
		seasons = append(seasons, season.ToModel())
	}
	return seasons, nil
}

func (s *projectService) validateSeasonUniqueness(ctx context.Context, showID string, title string) error {
	seasons, err := s.repository.ListShowSeasons(ctx, showID)
	if err != nil {
		return errors.Wrap(err, "Error while getting seasons from the Mongo database")
	}
	for _, season := range seasons {
		if season.Title == title {
			return errors.New("There is already a season with that name in the Mongo database with id:" + season.ID)
		}
	}
	return nil
}

func toSeasonDTO(ID string, showID string, title string, trailerURL string, postersPath []string, releaseDate time.Time, rating float64, resume string, directedBy models.FilmCrews, producedBy models.FilmCrews, writtenBy models.FilmCrews, episodes models.ShortEpisodes) *dto.SeasonDTO {
	return &dto.SeasonDTO{
		ID:          ID,
		ShowID:      showID,
		Title:       title,
		TrailerURL:  trailerURL,
		PostersPath: postersPath,
		Resume:      resume,
		Rating:      rating,
		ReleaseDate: releaseDate,
		WrittenBy:   toFilmCrewsDTO(writtenBy),
		ProducedBy:  toFilmCrewsDTO(producedBy),
		DirectedBy:  toFilmCrewsDTO(directedBy),
		Episodes:    toShortEpisodesDTO(episodes),
	}
}

func toShortSeasonDTO(seasonsModel models.ShortSeasons) dto.ShortSeasonsDTO {
	shortSeasons := dto.ShortSeasonsDTO{}
	for _, season := range seasonsModel {
		shortSeasons = append(shortSeasons, &dto.ShortSeasonDTO{
			ID:          season.ID,
			Title:       season.Title,
			PostersPath: season.PostersPath,
			Rating:      season.Rating,
		})
	}
	return shortSeasons
}
