package service

import (
	"context"
	"int-service/dto"
	"int-service/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type JournalistServicer interface {
	CreateJournalist(ctx context.Context, name string) (models.ResponseModeler, error)
	GetJournalist(ctx context.Context, ID string) (models.ResponseModeler, error)
	UpdateJournalist(ctx context.Context, ID string, name string) (models.ResponseModeler, error)
	ListJournalists(ctx context.Context) ([]models.ResponseModeler, error)
	GetJournalistByName(ctx context.Context, name string) (models.ResponseModeler, error)
}

func (s *projectService) CreateJournalist(ctx context.Context, name string) (models.ResponseModeler, error) {
	_, err := s.repository.GetJournalistByName(ctx, name)
	if err == nil {
		return nil, errors.New("There is already a journalist with that name.")
	}

	journalist := dto.JournalistDTO{
		ID:   uuid.New().String(),
		Name: name,
	}
	resp, err := s.repository.CreateJournalist(ctx, &journalist)
	if err != nil {
		s.logger.Error("Error while creating journalist")
		return nil, errors.Wrap(err, "Error while creating journalist")
	}
	return resp.ToModel(), nil
}
func (s *projectService) GetJournalistByName(ctx context.Context, name string) (models.ResponseModeler, error) {
	resp, err := s.repository.GetJournalistByName(ctx, name)
	if err != nil {
		s.logger.Error("Error while getting journalist by name")
		return nil, errors.Wrap(err, "Error while getting journalist by name")
	}
	return resp.ToModel(), nil
}

func (s *projectService) GetJournalist(ctx context.Context, ID string) (models.ResponseModeler, error) {
	resp, err := s.repository.GetJournalist(ctx, ID)
	if err != nil {
		s.logger.Error("Error while getting journalist by id")
		return nil, errors.Wrap(err, "Error while getting journalist by id")
	}
	return resp.ToModel(), nil
}

func (s *projectService) UpdateJournalist(ctx context.Context, ID string, name string) (models.ResponseModeler, error) {
	journalist := dto.JournalistDTO{
		ID:   ID,
		Name: name,
	}
	resp, err := s.repository.UpdateJournalist(ctx, &journalist)
	if err != nil {
		s.logger.Error("Error while updating journalist")
		return nil, errors.New("Error while updating journalist")
	}
	return resp.ToModel(), nil
}

func (s *projectService) ListJournalists(ctx context.Context) ([]models.ResponseModeler, error) {
	resp, err := s.repository.ListJournalists(ctx)
	if err != nil {
		s.logger.Error("Error while listing all journalists")
		return nil, errors.New("Error while listing all journalists")
	}
	journalists := []models.ResponseModeler{}
	for _, j := range resp {
		journalists = append(journalists, j.ToModel())
	}
	return journalists, nil
}
