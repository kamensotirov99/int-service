package repository

import (
	"context"
	"int-service/dto"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type GenreRepository interface {
	CreateGenre(ctx context.Context, newGenre *dto.GenreDTO) (*dto.GenreDTO, error)
	GetGenre(ctx context.Context, ID string) (*dto.GenreDTO, error)
	GetGenreByName(ctx context.Context, name string)(*dto.GenreDTO,error)
	UpdateGenre(ctx context.Context, updatedGenre *dto.GenreDTO) (*dto.GenreDTO, error)
	ListGenres(ctx context.Context) (dto.GenresDTO, error)
}

func (m *MongoDatabase) CreateGenre(ctx context.Context, newGenre *dto.GenreDTO) (*dto.GenreDTO, error) {
	collection := m.client.Database("Project").Collection("Genres")
	_, err := collection.InsertOne(ctx, newGenre)
	if err != nil {
		return nil, errors.Wrap(err, "Error while inserting the new genre in the Mongo database")
	}
	return newGenre, nil
}

func (m *MongoDatabase) GetGenreByName(ctx context.Context, name string) (*dto.GenreDTO, error) {
	collection := m.client.Database("Project").Collection("Genres")
	filter := bson.D{bson.E{Key: "name", Value: name}}
	genre := dto.GenreDTO{}

	err := collection.FindOne(ctx, filter).Decode(&genre)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding genre by name from the Mongo database")
	}
	return &genre, nil
}

func (m *MongoDatabase) GetGenre(ctx context.Context, ID string) (*dto.GenreDTO, error) {
	collection := m.client.Database("Project").Collection("Genres")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	genre := dto.GenreDTO{}

	err := collection.FindOne(ctx, filter).Decode(&genre)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding genre by id from the Mongo database")
	}
	return &genre, nil
}

func (m *MongoDatabase) UpdateGenre(ctx context.Context, updatedGenre *dto.GenreDTO) (*dto.GenreDTO, error) {
	collection := m.client.Database("Project").Collection("Genres")
	filter := bson.D{bson.E{Key: "id", Value: updatedGenre.ID}}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "name", Value: updatedGenre.Name},
			bson.E{Key: "description", Value: updatedGenre.Description},
		}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating genre in the Mongo database")
	}
	return updatedGenre, nil
}

func (m *MongoDatabase) ListGenres(ctx context.Context) (dto.GenresDTO, error) {
	collection := m.client.Database("Project").Collection("Genres")
	genres := dto.GenresDTO{}

	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding all genres from the Mongo database")
	}

	for cursor.Next(ctx) {
		genre := dto.GenreDTO{}
		err = cursor.Decode(&genre)
		if err != nil {
			return nil, errors.Wrap(err, "Error while decoding genre")
		}
		genres = append(genres, &genre)
	}
	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error with the cursor")
	}
	cursor.Close(ctx)

	return genres, nil
}
