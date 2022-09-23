package service

import (
	"context"
	"int-service/dto"
	"int-service/models"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ArticleServicer interface {
	CreateArticle(ctx context.Context, title string, releaseDate time.Time, postersPath []string, description string, journalistName string) (models.ResponseModeler, error)
	GetArticle(ctx context.Context, ID string) (models.ResponseModeler, error)
	UpdateArticle(ctx context.Context, ID string, title string, releaseDate time.Time, postersPath []string, description string, journalistModel *models.Journalist) (models.ResponseModeler, error)
	ListArticles(ctx context.Context,elementCount int) ([]models.ResponseModeler, error)
	ListArticlesByJournalist(ctx context.Context, journalistID string) ([]models.ResponseModeler, error)
	UploadArticlePosters(ctx context.Context, ID string, postersPath []string) (models.ResponseModeler, error)
	DeleteArticlePoster(ctx context.Context, ID string, image string) error
}

func (s *projectService) CreateArticle(ctx context.Context, title string, releaseDate time.Time, postersPath []string, description string, journalistName string) (models.ResponseModeler, error) {
	search, err := s.repository.GetJournalistByName(ctx, journalistName)
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting journalist with name : "+search.Name)
	}
	article := toArticleDTO(uuid.New().String(), title, releaseDate, postersPath, description, search.ID)
	resp, err := s.repository.CreateArticle(ctx, article)
	if err != nil {
		s.logger.Error("Error while creating article")
		return nil, errors.Wrap(err, "Error while creating article")
	}
	return resp.ToModel(), nil
}

func (s *projectService) GetArticle(ctx context.Context, ID string) (models.ResponseModeler, error) {
	resp, err := s.repository.GetArticle(ctx, ID)
	if err != nil {
		s.logger.Error("Error while getting article by id")
		return nil, errors.Wrap(err, "Error while getting article by id")
	}
	return resp.ToModel(), nil
}

func (s *projectService) UpdateArticle(ctx context.Context, ID string, title string, releaseDate time.Time, postersPath []string, description string, journalistModel *models.Journalist) (models.ResponseModeler, error) {
	updatedArticle := toArticleDTO(ID, title, releaseDate, postersPath, description, journalistModel.ID)
	resp, err := s.repository.UpdateArticle(ctx, updatedArticle)
	if err != nil {
		s.logger.Error("Error while updating article")
		return nil, errors.New("Error while updating article")
	}
	return resp.ToModel(), nil
}

func (s *projectService) ListArticles(ctx context.Context,elementCount int) ([]models.ResponseModeler, error) {
	resp, err := s.repository.ListArticles(ctx,elementCount)
	if err != nil {
		s.logger.Error("Error while listing all articles")
		return nil, errors.New("Error while listing all articles")
	}
	articles := []models.ResponseModeler{}
	for _, article := range resp {
		articles = append(articles, article.ToModel())
	}
	return articles, nil
}

func (s *projectService) ListArticlesByJournalist(ctx context.Context, journalistID string) ([]models.ResponseModeler, error) {
	resp, err := s.repository.ListArticlesByJournalist(ctx, journalistID)
	if err != nil {
		s.logger.Error("Error while listing all articles by journalist Id")
		return nil, errors.New("Error while listing all articles by journalist Id")
	}
	articles := []models.ResponseModeler{}
	for _, article := range resp {
		articles = append(articles, article.ToModel())
	}
	return articles, nil
}

func (s *projectService) UploadArticlePosters(ctx context.Context, ID string, postersPath []string) (models.ResponseModeler, error) {
	resp, err := s.repository.UploadArticlePosters(ctx, ID, postersPath)
	if err != nil {
		s.logger.Error("Error while updating article posters")
		return nil, errors.Wrap(err, "Error while updating article posters")
	}
	return resp.ToModel(), nil
}

func (s *projectService) DeleteArticlePoster(ctx context.Context, ID string, image string) error {
	err := s.repository.DeleteArticlePoster(ctx, ID, image)
	if err != nil {
		s.logger.Error("Error while updating deleted article poster in database")
		return errors.Wrap(err, "Error while updating deleted article poster in database")
	}
	return nil
}

func toArticleDTO(ID string, title string, releaseDate time.Time, postersPath []string, description string, journalistID string) *dto.ArticleDTO {
	journalist := dto.ShortJournalistDTO{
		ID: journalistID,
	}
	return &dto.ArticleDTO{
		ID:          ID,
		Title:       title,
		ReleaseDate: releaseDate,
		PostersPath: postersPath,
		Description: description,
		Journalist:  journalist,
	}
}
