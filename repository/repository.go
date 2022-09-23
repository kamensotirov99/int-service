package repository

import (
	"context"
	"int-service/dto"
)

type Repository interface {
	CreateClothing(ctx context.Context, clothing *dto.ClothingDTO) (*dto.ClothingDTO, error)
	DeleteClothing(ctx context.Context, ID string) error
	GetAll(ctx context.Context) (*dto.ClothesDTO, error)
}

type ProjectRepository interface {
	ShowRepository
	SeasonRepository
	CelebrityRepository
	EpisodeRepository
	ArticleRepository
	GenreRepository
	JournalistRepository
}