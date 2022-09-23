package service

import (
	"context"
	"int-service/dto"
	"int-service/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type GenreServicer interface {
	CreateGenre(ctx context.Context, name string, description string) (models.ResponseModeler, error)
	GetGenre(ctx context.Context, ID string) (models.ResponseModeler, error)
	GetGenreByName(ctx context.Context, name string) (models.ResponseModeler, error)
	UpdateGenre(ctx context.Context, ID string, name string, description string) (models.ResponseModeler, error)
	ListGenres(ctx context.Context) ([]models.ResponseModeler, error)
}

func (s *projectService) CreateGenre(ctx context.Context, name string, description string) (models.ResponseModeler, error) {
	_, err := s.repository.GetGenreByName(ctx, name)
	if err == nil {
		return nil, errors.New("There is already a genre with that name.")
	}

	genre := dto.GenreDTO{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
	}
	resp, err := s.repository.CreateGenre(ctx, &genre)
	if err != nil {
		s.logger.Error("Error while creating genre")
		return nil, errors.Wrap(err, "Error while creating genre")
	}
	return resp.ToModel(), nil
}

func (s *projectService) GetGenreByName(ctx context.Context, name string) (models.ResponseModeler, error) {
	resp, err := s.repository.GetGenreByName(ctx, name)
	if err != nil {
		s.logger.Error("Error while getting genre by name")
		return nil, errors.Wrap(err, "Error while getting genre by name")
	}
	return resp.ToModel(), nil
}

func (s *projectService) GetGenre(ctx context.Context, ID string) (models.ResponseModeler, error) {
	resp, err := s.repository.GetGenre(ctx, ID)
	if err != nil {
		s.logger.Error("Error while getting genre by id")
		return nil, errors.Wrap(err, "Error while getting genre by id")
	}
	return resp.ToModel(), nil
}

func (s *projectService) UpdateGenre(ctx context.Context, ID string, name string, description string) (models.ResponseModeler, error) {
	updatedGenre := &dto.GenreDTO{
		ID:          ID,
		Name:        name,
		Description: description,
	}
	resp, err := s.repository.UpdateGenre(ctx, updatedGenre)
	if err != nil {
		s.logger.Error("Error while updating genre")
		return nil, errors.New("Error while updating genre")
	}
	return resp.ToModel(), nil
}

func (s *projectService) ListGenres(ctx context.Context) ([]models.ResponseModeler, error) {
	resp, err := s.repository.ListGenres(ctx)
	if err != nil {
		s.logger.Error("Error while listing all genres")
		return nil, errors.New("Error while listing all genres")
	}
	genres := []models.ResponseModeler{}
	for _, genre := range resp {
		genres = append(genres, genre.ToModel())
	}
	return genres, nil
}

func toShortGenresDTO(genresModel models.ShortGenres) dto.ShortGenresDTO {
	genres := dto.ShortGenresDTO{}
	for _, genre := range genresModel {
		genres = append(genres, &dto.ShortGenreDTO{
			ID:          genre.ID,
			Name:        genre.Name,
		})
	}
	return genres
}
