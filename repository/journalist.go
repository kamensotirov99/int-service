package repository

import (
	"context"
	"int-service/dto"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type JournalistRepository interface {
	CreateJournalist(ctx context.Context, newJournalist *dto.JournalistDTO) (*dto.JournalistDTO, error)
	GetJournalist(ctx context.Context, ID string) (*dto.JournalistDTO, error)
	UpdateJournalist(ctx context.Context, updatedJournalist *dto.JournalistDTO) (*dto.JournalistDTO, error)
	ListJournalists(ctx context.Context) (dto.JournalistsDTO, error)
	GetJournalistByName(ctx context.Context, name string) (*dto.JournalistDTO, error)
}

func (m *MongoDatabase) CreateJournalist(ctx context.Context, newJournalist *dto.JournalistDTO) (*dto.JournalistDTO, error) {
	collection := m.client.Database("Project").Collection("Journalists")
	_, err := collection.InsertOne(ctx, newJournalist)
	if err != nil {
		return nil, errors.Wrap(err, "Error while inserting the new journalist in the Mongo database")
	}
	return newJournalist, nil
}

func (m *MongoDatabase) GetJournalistByName(ctx context.Context, name string) (*dto.JournalistDTO, error) {
	collection := m.client.Database("Project").Collection("Journalists")
	filter := bson.D{bson.E{Key: "name", Value: name}}
	journalist := dto.JournalistDTO{}

	err := collection.FindOne(ctx, filter).Decode(&journalist)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding journalist by name from the Mongo database")
	}
	return &journalist, nil
}

func (m *MongoDatabase) GetJournalist(ctx context.Context, ID string) (*dto.JournalistDTO, error) {
	collection := m.client.Database("Project").Collection("Journalists")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	journalist := dto.JournalistDTO{}

	err := collection.FindOne(ctx, filter).Decode(&journalist)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding journalist by id from the Mongo database")
	}
	return &journalist, nil
}

func (m *MongoDatabase) UpdateJournalist(ctx context.Context, updatedJournalist *dto.JournalistDTO) (*dto.JournalistDTO, error) {
	collection := m.client.Database("Project").Collection("Journalists")
	filter := bson.D{bson.E{Key: "id", Value: updatedJournalist.ID}}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "name", Value: updatedJournalist.Name},
		}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating journalist in the Mongo database")
	}
	return updatedJournalist, nil
}

func (m *MongoDatabase) ListJournalists(ctx context.Context) (dto.JournalistsDTO, error) {
	collection := m.client.Database("Project").Collection("Journalists")
	journalists := dto.JournalistsDTO{}

	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding all journalists from the Mongo database")
	}

	for cursor.Next(ctx) {
		journalist := dto.JournalistDTO{}
		err = cursor.Decode(&journalist)
		if err != nil {
			return nil, errors.Wrap(err, "Error while decoding journalist")
		}
		journalists = append(journalists, &journalist)
	}
	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error with the cursor")
	}
	cursor.Close(ctx)

	return journalists, nil
}
