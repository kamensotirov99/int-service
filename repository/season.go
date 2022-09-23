package repository

import (
	"context"
	"int-service/dto"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SeasonRepository interface {
	CreateSeason(ctx context.Context, newSeason *dto.SeasonDTO) (*dto.SeasonDTO, error)
	AddShortEpisode(ctx context.Context, seasonID string, newEpisode *dto.ShortEpisodeDTO) (*dto.ShortEpisodeDTO, error)
	GetSeason(ctx context.Context, seasonID string) (*dto.SeasonDTO, error)
	UpdateSeason(ctx context.Context, updatedSeason *dto.SeasonDTO) (*dto.SeasonDTO, error)
	UpdateShortEpisode(ctx context.Context, updatedEpisode *dto.ShortEpisodeDTO) (*dto.ShortEpisodeDTO, error)
	UpdateShortCelebritiesInSeasons(ctx context.Context, updatedCelebrity *dto.ShortCelebrityDTO, celebrityType string) (*dto.ShortCelebrityDTO, error)
	UploadSeasonPosters(ctx context.Context, seasonID string, postersPath []string) (*dto.SeasonDTO, error)
	DeleteSeasonPoster(ctx context.Context, seriesID string, seasonID string, image string) error
	DeleteShortCelebritiesPostersInSeason(ctx context.Context, celebrityID string, image string, celebrityType string) error
	ListShowSeasons(ctx context.Context, ID string) (dto.SeasonsDTO, error)
	ListSeasonsCollection(ctx context.Context) (dto.SeasonsDTO, error)
}

func (m *MongoDatabase) CreateSeason(ctx context.Context, newSeason *dto.SeasonDTO) (*dto.SeasonDTO, error) {
	collection := m.client.Database("Project").Collection("Seasons")
	newSeason.PostersPath = []string{}
	_, err := collection.InsertOne(ctx, newSeason)
	if err != nil {
		return nil, errors.Wrap(err, "Error while inserting new season in the Mongo database")
	}
	return newSeason, nil
}

func (m *MongoDatabase) AddShortEpisode(ctx context.Context, seasonID string, newEpisode *dto.ShortEpisodeDTO) (*dto.ShortEpisodeDTO, error) {
	collection := m.client.Database("Project").Collection("Seasons")
	filter := bson.D{bson.E{Key: "id", Value: seasonID}}
	update := bson.M{"$push": bson.M{
		"episodes": bson.M{
			"id":          newEpisode.ID,
			"title":       newEpisode.Title,
			"postersPath": newEpisode.PostersPath,
			"rating":      newEpisode.Rating,
			"resume":      newEpisode.Resume,
		}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errors.Wrap(err, "Error while adding short episode in the Mongo database")
	}
	return &dto.ShortEpisodeDTO{
		ID:          newEpisode.ID,
		Title:       newEpisode.Title,
		PostersPath: newEpisode.PostersPath,
		Rating:      newEpisode.Rating,
		Resume:      newEpisode.Resume,
	}, nil
}

func (m *MongoDatabase) GetSeason(ctx context.Context, ID string) (*dto.SeasonDTO, error) {
	collection := m.client.Database("Project").Collection("Seasons")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	season := dto.SeasonDTO{}

	err := collection.FindOne(ctx, filter).Decode(&season)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding season by id from the Mongo database")
	}
	return &season, nil
}

func (m *MongoDatabase) UpdateSeason(ctx context.Context, updatedSeason *dto.SeasonDTO) (*dto.SeasonDTO, error) {
	collection := m.client.Database("Project").Collection("Seasons")
	filter := bson.D{bson.E{Key: "id", Value: updatedSeason.ID}}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "title", Value: updatedSeason.Title},
			bson.E{Key: "trailerUrl", Value: updatedSeason.TrailerURL},
			bson.E{Key: "postersPath", Value: updatedSeason.PostersPath},
			bson.E{Key: "resume", Value: updatedSeason.Resume},
			bson.E{Key: "rating", Value: updatedSeason.Rating},
			bson.E{Key: "releaseDate", Value: updatedSeason.ReleaseDate},
			bson.E{Key: "writtenBy", Value: updatedSeason.WrittenBy},
			bson.E{Key: "producedBy", Value: updatedSeason.ProducedBy},
			bson.E{Key: "directedBy", Value: updatedSeason.DirectedBy},
			bson.E{Key: "episodes", Value: updatedSeason.Episodes},
		}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating season in the Mongo database")
	}
	return updatedSeason, nil
}

func (m *MongoDatabase) UpdateShortEpisode(ctx context.Context, updatedEpisode *dto.ShortEpisodeDTO) (*dto.ShortEpisodeDTO, error) {
	collection := m.client.Database("Project").Collection("Seasons")
	_, err := collection.UpdateOne(
		ctx,
		bson.D{bson.E{Key: "episodes.id", Value: updatedEpisode.ID}},
		bson.M{"$set": bson.M{
			"episodes.$[elem].title":       updatedEpisode.Title,
			"episodes.$[elem].postersPath": updatedEpisode.PostersPath,
			"episodes.$[elem].rating":      updatedEpisode.Rating,
			"episodes.$[elem].resume":      updatedEpisode.Resume,
		}},

		options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{bson.M{"elem.id": updatedEpisode.ID}},
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating short episode in the Mongo database")
	}
	return updatedEpisode, nil
}

func (m *MongoDatabase) UploadSeasonPosters(ctx context.Context, seasonID string, postersPath []string) (*dto.SeasonDTO, error) {
	collection := m.client.Database("Project").Collection("Seasons")
	filter := bson.D{bson.E{Key: "id", Value: seasonID}}
	update := bson.M{"$push": bson.M{
		"postersPath": bson.M{
			"$each": postersPath}}}

	updatedSeason := dto.SeasonDTO{}

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	resp := collection.FindOneAndUpdate(ctx, filter, update, &opt)
	err := resp.Decode(&updatedSeason)
	if err != nil {
		return nil, errors.Wrap(err, "Error while uploading season posters in the Mongo database")
	}

	return &updatedSeason, nil
}

func (m *MongoDatabase) DeleteSeasonPoster(ctx context.Context, seriesID string, seasonID string, image string) error {
	collection := m.client.Database("Project").Collection("Seasons")
	filter := bson.D{bson.E{Key: "id", Value: seasonID}}
	posterPath := "/series/" + seriesID + "/" + seasonID + "/" + image
	update := bson.M{"$pull": bson.M{"postersPath": posterPath}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "Error while deleting season poster in the Mongo database")
	}

	return nil
}

func (m *MongoDatabase) DeleteShortCelebritiesPostersInSeason(ctx context.Context, celebrityID string, image string, celebrityType string) error {
	collection := m.client.Database("Project").Collection("Seasons")
	filter := bson.D{bson.E{Key: celebrityType + ".id", Value: celebrityID}}
	posterPath := "/celebrities/" + celebrityID + "/" + image
	update := bson.M{"$pull": bson.M{celebrityType + ".$.postersPath": posterPath}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "Error while deleting celebrity poster in short celebrity in the Mongo database")
	}
	return nil
}

func (m *MongoDatabase) ListShowSeasons(ctx context.Context, showID string) (dto.SeasonsDTO, error) {
	collection := m.client.Database("Project").Collection("Seasons")
	filter := bson.D{bson.E{Key: "showId", Value: showID}}
	seasons := dto.SeasonsDTO{}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding all seasons from the Mongo database")
	}
	for cursor.Next(ctx) {
		season := dto.SeasonDTO{}
		err = cursor.Decode(&season)
		if err != nil {
			return nil, errors.Wrap(err, "Error while decoding showSeason")
		}
		seasons = append(seasons, &season)
	}
	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error with the cursor")
	}
	cursor.Close(ctx)

	return seasons, nil
}

func (m *MongoDatabase) ListSeasonsCollection(ctx context.Context) (dto.SeasonsDTO, error) {
	collection := m.client.Database("Project").Collection("Seasons")
	seasons := dto.SeasonsDTO{}

	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding all show seasons from the Mongo database")
	}

	for cursor.Next(ctx) {
		season := dto.SeasonDTO{}
		err = cursor.Decode(&season)
		if err != nil {
			return nil, errors.Wrap(err, "Error while decoding showSeason")
		}
		seasons = append(seasons, &season)
	}
	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error with the cursor")
	}
	cursor.Close(ctx)

	return seasons, nil
}

func (m *MongoDatabase) UpdateShortCelebritiesInSeasons(ctx context.Context, updatedCelebrity *dto.ShortCelebrityDTO, celebrityType string) (*dto.ShortCelebrityDTO, error) {
	collection := m.client.Database("Project").Collection("Seasons")
	_, err := collection.UpdateMany(
		ctx,
		bson.D{bson.E{Key: celebrityType + ".id", Value: updatedCelebrity.ID}},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: celebrityType + ".$.name", Value: updatedCelebrity.Name}}},
			{Key: "$push", Value: bson.D{
				{Key: celebrityType + ".$.postersPath", Value: bson.D{
					{Key: "$each", Value: updatedCelebrity.PostersPath}}}}},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating short celebrity in the Mongo database")
	}
	return updatedCelebrity, nil
}
