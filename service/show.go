package service

import (
	"context"
	"int-service/dto"
	"int-service/models"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ShowServicer interface {
	CreateShow(ctx context.Context, title string, sType string, postersPath []string, releaseDate time.Time, endDate time.Time, rating float64, length *models.ShowLength, trailerURL string, genres models.ShortGenres, directedBy models.FilmCrews, producedBy models.FilmCrews, writtenBy models.FilmCrews, starring models.ShortCelebrities, description string, seasons models.ShortSeasons) (models.ResponseModeler, error)
	GetShow(ctx context.Context, ID string) (models.ResponseModeler, error)
	UpdateShow(ctx context.Context, ID string, title string, sType string, postersPath []string, releaseDate time.Time, endDate time.Time, rating float64, length *models.ShowLength, trailerURL string, genres models.ShortGenres, directedBy models.FilmCrews, producedBy models.FilmCrews, writtenBy models.FilmCrews, starring models.ShortCelebrities, description string, seasons models.ShortSeasons) (models.ResponseModeler, error)
	ListShows(ctx context.Context) ([]models.ResponseModeler, error)
	UploadSeriesPosters(ctx context.Context, ID string, postersPath []string) (models.ResponseModeler, error)
	DeleteSeriesPoster(ctx context.Context, ID string, image string) error
	UploadMoviePosters(ctx context.Context, ID string, postersPath []string) (models.ResponseModeler, error)
	DeleteMoviePoster(ctx context.Context, ID string, image string) error
}

func (s *projectService) CreateShow(ctx context.Context, title string, sType string, postersPath []string, releaseDate time.Time, endDate time.Time, rating float64, length *models.ShowLength, trailerURL string, genres models.ShortGenres, directedBy models.FilmCrews, producedBy models.FilmCrews, writtenBy models.FilmCrews, starring models.ShortCelebrities, description string, seasons models.ShortSeasons) (models.ResponseModeler, error) {
	err := s.validateShowUniqueness(ctx, title, releaseDate)
	if err != nil {
		s.logger.Error("Error while creating show")
		return nil, errors.Wrap(err, "Error while creating show")
	}
	show := toShowDTO(uuid.New().String(), title, sType, postersPath, releaseDate, endDate, rating, length, trailerURL, genres, directedBy, producedBy, writtenBy, starring, description, seasons)
	resp, err := s.repository.CreateShow(ctx, show)
	if err != nil {
		s.logger.Error("Error while creating show")
		return nil, errors.Wrap(err, "Error while creating show")
	}
	return resp.ToModel(), nil
}

func (s *projectService) GetShow(ctx context.Context, ID string) (models.ResponseModeler, error) {
	resp, err := s.repository.GetShow(ctx, ID)
	if err != nil {
		s.logger.Error("Error while getting show by id")
		return nil, errors.Wrap(err, "Error while getting show by id")
	}
	return resp.ToModel(), nil
}

func (s *projectService) UpdateShow(ctx context.Context, ID string, title string, sType string, postersPath []string, releaseDate time.Time, endDate time.Time, rating float64, length *models.ShowLength, trailerURL string, genres models.ShortGenres, directedBy models.FilmCrews, producedBy models.FilmCrews, writtenBy models.FilmCrews, starring models.ShortCelebrities, description string, seasons models.ShortSeasons) (models.ResponseModeler, error) {
	updatedShow := toShowDTO(ID, title, sType, postersPath, releaseDate, endDate, rating, length, trailerURL, genres, directedBy, producedBy, writtenBy, starring, description, seasons)
	resp, err := s.repository.UpdateShow(ctx, updatedShow)
	if err != nil {
		s.logger.Error("Error while updating show")
		return nil, errors.New("Error while updating show")
	}
	return resp.ToModel(), nil
}

func (s *projectService) ListShows(ctx context.Context) ([]models.ResponseModeler, error) {
	resp, err := s.repository.ListShows(ctx)
	if err != nil {
		s.logger.Error("Error while listing all shows")
		return nil, errors.New("Error while listing all shows")
	}
	shows := []models.ResponseModeler{}
	for _, show := range resp {
		shows = append(shows, show.ToModel())
	}
	return shows, nil
}

func (s *projectService) UploadSeriesPosters(ctx context.Context, ID string, postersPath []string) (models.ResponseModeler, error) {
	resp, err := s.repository.UploadSeriesPosters(ctx, ID, postersPath)
	if err != nil {
		s.logger.Error("Error while updating series posters")
		return nil, errors.Wrap(err, "Error while updating series posters")
	}
	return resp.ToModel(), nil
}

func (s *projectService) DeleteSeriesPoster(ctx context.Context, ID string, image string) error {
	err := s.repository.DeleteSeriesPoster(ctx, ID, image)
	if err != nil {
		s.logger.Error("Error while updating deleted series poster in database")
		return errors.Wrap(err, "Error while updating deleted series poster in database")
	}
	return nil
}

func (s *projectService) UploadMoviePosters(ctx context.Context, ID string, postersPath []string) (models.ResponseModeler, error) {
	resp, err := s.repository.UploadMoviePosters(ctx, ID, postersPath)
	if err != nil {
		s.logger.Error("Error while updating movie posters")
		return nil, errors.Wrap(err, "Error while updating movie posters")
	}
	return resp.ToModel(), nil
}

func (s *projectService) DeleteMoviePoster(ctx context.Context, ID string, image string) error {
	err := s.repository.DeleteMoviePoster(ctx, ID, image)
	if err != nil {
		s.logger.Error("Error while deleting movie poster in database")
		return errors.Wrap(err, "Error while deleting movie poster in database")
	}
	return nil
}

func (s *projectService) validateShowUniqueness(ctx context.Context, title string, releaseDate time.Time) error {
	shows, err := s.repository.ListShows(ctx)
	if err != nil {
		return errors.Wrap(err, "Error while getting shows from the Mongo database")
	}
	for _, show := range shows {
		if show.Title == title && show.ReleaseDate == releaseDate {
			return errors.New("There is already a show with that name and release date.")
		}
	}
	return nil
}

func toShowDTO(ID string, title string, sType string, postersPath []string, releaseDate time.Time, endDate time.Time, rating float64, lengthModel *models.ShowLength, trailerURL string, genres models.ShortGenres, directedBy models.FilmCrews, producedBy models.FilmCrews, writtenBy models.FilmCrews, starring models.ShortCelebrities, description string, seasons models.ShortSeasons) *dto.ShowDTO {
	length := dto.ShowLengthDTO{
		Hours:   lengthModel.Hours,
		Minutes: lengthModel.Minutes,
	}
	show := &dto.ShowDTO{
		ID:          ID,
		Title:       title,
		Type:        sType,
		PostersPath: postersPath,
		ReleaseDate: releaseDate,
		EndDate:     endDate,
		Rating:      rating,
		Length:      length,
		TrailerURL:  trailerURL,
		Genres:      toShortGenresDTO(genres),
		WrittenBy:   toFilmCrewsDTO(writtenBy),
		ProducedBy:  toFilmCrewsDTO(producedBy),
		DirectedBy:  toFilmCrewsDTO(directedBy),
		Starring:    toShortCelebDTO(starring),
		Description: description,
		Seasons:     toShortSeasonDTO(seasons),
	}
	return show
}
