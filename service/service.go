package service

import (
	"context"
	"int-service/dto"
	"int-service/models"
	"int-service/repository"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type service struct {
	logger     *logrus.Logger
	repository repository.Repository
}

type Servicer interface {
	CreateClothing(ctx context.Context, Type string, Size int, Price int, Gender string) (models.ResponseModeler, error)
	DeleteClothing(ctx context.Context, ID string) error
	GetAll(ctx context.Context) ([]models.ResponseModeler, error)
}

type projectService struct {
	logger     *logrus.Logger
	repository repository.ProjectRepository
}

type ProjectServicer interface {
	ArticleServicer
	CelebrityServicer
	EpisodeServicer
	ShowServicer
	GenreServicer
	SeasonServicer
	JournalistServicer
}

func New(logger *logrus.Logger, repository repository.Repository) Servicer {
	return &service{logger, repository}
}

func NewSvc(logger *logrus.Logger, repository repository.ProjectRepository) ProjectServicer {
	return &projectService{logger, repository}
}

func (s *service) CreateClothing(ctx context.Context, cType string, size int, price int, gender string) (models.ResponseModeler, error) {
	c := dto.ClothingDTO{
		ID:     uuid.New().String(),
		Type:   cType,
		Size:   size,
		Price:  price,
		Gender: gender,
	}
	resp, err := s.repository.CreateClothing(ctx, &c)
	if err != nil {
		s.logger.Error("Error while creating clothing")
		return nil, errors.Wrap(err, "Error while creating clothing")
	}

	return resp.ToModel(), nil
}

func (s *service) DeleteClothing(ctx context.Context, ID string) error {
	err := s.repository.DeleteClothing(ctx, ID)
	if err != nil {
		s.logger.Error("Error while deleting clothing")
		return errors.Wrap(err, "Error while deleting clothing")
	}
	return nil
}

func (s *service) GetAll(ctx context.Context) ([]models.ResponseModeler, error) {
	resp, err := s.repository.GetAll(ctx)
	if err != nil {
		s.logger.Error("Error while gettilg all clothes")
		return nil, errors.Wrap(err, "Error while getting all clothes")
	}

	clothes := []models.ResponseModeler{}
	for _, c := range resp.Clothes {
		clothes = append(clothes, c.ToModel())
	}
	return clothes, nil
}
