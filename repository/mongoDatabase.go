package repository

import (
	"context"
	"int-service/dto"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDatabase struct {
	client *mongo.Client
}

func NewMongoDatabase(c *mongo.Client) Repository {
	return &MongoDatabase{
		client: c,
	}
}

func NewMongoDB(c *mongo.Client) ProjectRepository {
	return &MongoDatabase{
		client: c,
	}
}

func (m *MongoDatabase) CreateClothing(ctx context.Context, newClothing *dto.ClothingDTO) (*dto.ClothingDTO, error) {
	collection := m.client.Database("Clothing").Collection("Summer")
	_, err := collection.InsertOne(ctx, newClothing)
	if err != nil {
		return nil, errors.Wrap(err, "Error while inserting the new clothing in the Mongo database")
	}

	return newClothing, nil
}

func (m *MongoDatabase) DeleteClothing(ctx context.Context, ID string) error {
	collection := m.client.Database("Clothing").Collection("Summer")
	filter := bson.D{bson.E{Key: "id", Value: ID}}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Error while deleting the clothing from the Mongo database")
	}
	return nil
}

func (m *MongoDatabase) GetAll(ctx context.Context) (*dto.ClothesDTO, error) {
	collection := m.client.Database("Clothing").Collection("Summer")
	clothes := dto.ClothesDTO{}

	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding clothes from the Mongo database")
	}

	for cursor.Next(ctx) {
		clothing := &dto.ClothingDTO{}
		err = cursor.Decode(&clothing)
		if err != nil {
			return nil, errors.Wrap(err, "Error while decoding the clothing")
		}
		clothes.Clothes = append(clothes.Clothes, *clothing)
	}
	if err = cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error while using the mongo cursor")
	}
	cursor.Close(ctx)

	return &clothes, nil
}
