package repository

import (
	"context"
	"int-service/dto"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CelebrityRepository interface {
	CreateCelebrity(ctx context.Context, newCelebrity *dto.CelebrityDTO) (*dto.CelebrityDTO, error)
	GetCelebrity(ctx context.Context, ID string) (*dto.CelebrityDTO, error)
	UpdateCelebrity(ctx context.Context, updatedCelebrity *dto.CelebrityDTO) (*dto.CelebrityDTO, error)
	UploadCelebrityPosters(ctx context.Context, ID string, postersPath []string) (*dto.CelebrityDTO, error)
	DeleteCelebrityPoster(ctx context.Context, ID string, image string) error
	ListCelebrities(ctx context.Context) (dto.CelebritiesDTO, error)
}

func (m *MongoDatabase) CreateCelebrity(ctx context.Context, newCelebrity *dto.CelebrityDTO) (*dto.CelebrityDTO, error) {
	collection := m.client.Database("Project").Collection("Celebrities")
	newCelebrity.PostersPath = []string{}
	_, err := collection.InsertOne(ctx, newCelebrity)
	if err != nil {
		return nil, errors.Wrap(err, "Error while inserting the new celebrity in the Mongo database")
	}
	return newCelebrity, nil
}

func (m *MongoDatabase) GetCelebrity(ctx context.Context, ID string) (*dto.CelebrityDTO, error) {
	collection := m.client.Database("Project").Collection("Celebrities")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	celebrity := dto.CelebrityDTO{}

	err := collection.FindOne(ctx, filter).Decode(&celebrity)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding celebrity by id from the Mongo database")
	}
	return &celebrity, nil
}

func (m *MongoDatabase) UpdateCelebrity(ctx context.Context, updatedCelebrity *dto.CelebrityDTO) (*dto.CelebrityDTO, error) {
	collection := m.client.Database("Project").Collection("Celebrities")
	filter := bson.D{bson.E{Key: "id", Value: updatedCelebrity.ID}}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "name", Value: updatedCelebrity.Name},
			bson.E{Key: "occupation", Value: updatedCelebrity.Occupation},
			bson.E{Key: "postersPath", Value: updatedCelebrity.PostersPath},
			bson.E{Key: "dateOfBirth", Value: updatedCelebrity.DateOfBirth},
			bson.E{Key: "dateOfDeath", Value: updatedCelebrity.DateOfDeath},
			bson.E{Key: "placeOfBirth", Value: updatedCelebrity.PlaceOfBirth},
			bson.E{Key: "gender", Value: updatedCelebrity.Gender},
			bson.E{Key: "bio", Value: updatedCelebrity.Bio},
		}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating celebrity in the Mongo database")
	}
	return updatedCelebrity, nil
}

func (m *MongoDatabase) UploadCelebrityPosters(ctx context.Context, ID string, postersPath []string) (*dto.CelebrityDTO, error) {
	collection := m.client.Database("Project").Collection("Celebrities")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	update := bson.M{"$push": bson.M{
		"postersPath": bson.M{
			"$each": postersPath}}}

	updatedCelebrity := dto.CelebrityDTO{}

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	resp := collection.FindOneAndUpdate(ctx, filter, update, &opt)
	err := resp.Decode(&updatedCelebrity)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating celebrity posters in the Mongo database")
	}

	return &updatedCelebrity, nil
}

func (m *MongoDatabase) DeleteCelebrityPoster(ctx context.Context, ID string, image string) error {
	collection := m.client.Database("Project").Collection("Celebrities")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	posterPath := "/celebrities/" + ID + "/" + image
	update := bson.M{"$pull": bson.M{"postersPath": posterPath}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "Error while updating deleted celebrity poster in the Mongo database")
	}

	return nil
}

func (m *MongoDatabase) ListCelebrities(ctx context.Context) (dto.CelebritiesDTO, error) {
	collection := m.client.Database("Project").Collection("Celebrities")
	celebrities := dto.CelebritiesDTO{}

	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding all celebrities from the Mongo database")
	}

	for cursor.Next(ctx) {
		celebrity := dto.CelebrityDTO{}
		err = cursor.Decode(&celebrity)
		if err != nil {
			return nil, errors.Wrap(err, "Error while decoding celebrity")
		}
		celebrities = append(celebrities, &celebrity)
	}
	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error with the cursor")
	}
	cursor.Close(ctx)

	return celebrities, nil
}
