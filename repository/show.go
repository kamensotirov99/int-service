package repository

import (
	"context"

	"int-service/dto"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShowRepository interface {
	CreateShow(ctx context.Context, newShow *dto.ShowDTO) (*dto.ShowDTO, error)
	AddShortSeason(ctx context.Context, showID string, season *dto.ShortSeasonDTO) (*dto.ShortSeasonDTO, error)
	GetShow(ctx context.Context, ID string) (*dto.ShowDTO, error)
	UpdateShow(ctx context.Context, updatedShow *dto.ShowDTO) (*dto.ShowDTO, error)
	UpdateShortSeason(ctx context.Context, updatedSeason *dto.ShortSeasonDTO) (*dto.ShortSeasonDTO, error)
	UpdateShortCelebritiesInShow(ctx context.Context, updatedCelebrity *dto.ShortCelebrityDTO, celebrityType string) (*dto.ShortCelebrityDTO, error)
	ListShows(ctx context.Context) (dto.ShowsDTO, error)
	UploadSeriesPosters(ctx context.Context, ID string, postersPath []string) (*dto.ShowDTO, error)
	DeleteSeriesPoster(ctx context.Context, ID string, image string) error
	UploadMoviePosters(ctx context.Context, ID string, postersPath []string) (*dto.ShowDTO, error)
	DeleteMoviePoster(ctx context.Context, ID string, image string) error
	DeleteShortCelebritiesPostersInShow(ctx context.Context, ID string, image string, celebrityType string) error
	DeleteShortSeasonPostersInShow(ctx context.Context, seriesID string, seasonID string, image string) error
}

func (m *MongoDatabase) CreateShow(ctx context.Context, newShow *dto.ShowDTO) (*dto.ShowDTO, error) {
	collection := m.client.Database("Project").Collection("Shows")
	newShow.PostersPath = []string{}
	_, err := collection.InsertOne(ctx, newShow)
	if err != nil {
		return nil, errors.Wrap(err, "Error while inserting the new show in the Mongo database")
	}
	return newShow, nil
}

func (m *MongoDatabase) AddShortSeason(ctx context.Context, showID string, newSeason *dto.ShortSeasonDTO) (*dto.ShortSeasonDTO, error) {
	collection := m.client.Database("Project").Collection("Shows")
	filter := bson.D{bson.E{Key: "id", Value: showID}}
	update := bson.M{"$push": bson.M{
		"seasons": bson.M{
			"id":          newSeason.ID,
			"title":       newSeason.Title,
			"postersPath": newSeason.PostersPath,
			"rating":      newSeason.Rating,
		}}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errors.Wrap(err, "Error while adding short season in the Mongo database")
	}
	return newSeason, nil
}

func (m *MongoDatabase) GetShow(ctx context.Context, ID string) (*dto.ShowDTO, error) {
	collection := m.client.Database("Project").Collection("Shows")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	show := dto.ShowDTO{}

	err := collection.FindOne(ctx, filter).Decode(&show)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding show by id from the Mongo database")
	}
	return &show, nil
}

func (m *MongoDatabase) UpdateShow(ctx context.Context, updatedShow *dto.ShowDTO) (*dto.ShowDTO, error) {
	collection := m.client.Database("Project").Collection("Shows")
	filter := bson.D{bson.E{Key: "id", Value: updatedShow.ID}}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "title", Value: updatedShow.Title},
			bson.E{Key: "releaseDate", Value: updatedShow.ReleaseDate},
			bson.E{Key: "type", Value: updatedShow.Type},
			bson.E{Key: "description", Value: updatedShow.Description},
			bson.E{Key: "endDate", Value: updatedShow.EndDate},
			bson.E{Key: "genres", Value: updatedShow.Genres},
			bson.E{Key: "postersPath", Value: updatedShow.PostersPath},
			bson.E{Key: "trailerUrl", Value: updatedShow.TrailerURL},
			bson.E{Key: "rating", Value: updatedShow.Rating},
			bson.E{Key: "length", Value: updatedShow.Length},
			bson.E{Key: "directedBy", Value: updatedShow.DirectedBy},
			bson.E{Key: "writtenBy", Value: updatedShow.WrittenBy},
			bson.E{Key: "producedBy", Value: updatedShow.ProducedBy},
			bson.E{Key: "starring", Value: updatedShow.Starring},
			bson.E{Key: "seasons", Value: updatedShow.Seasons},
		}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating show in the Mongo database")
	}
	return updatedShow, nil
}

func (m *MongoDatabase) UpdateShortSeason(ctx context.Context, updatedSeason *dto.ShortSeasonDTO) (*dto.ShortSeasonDTO, error) {
	collection := m.client.Database("Project").Collection("Shows")
	condition := bson.D{bson.E{Key: "seasons.id", Value: updatedSeason.ID}}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "seasons.$.title", Value: updatedSeason.Title},
			bson.E{Key: "seasons.$.rating", Value: updatedSeason.Rating},
			bson.E{Key: "seasons.$.postersPath", Value: updatedSeason.PostersPath},
		}},
	}
	_, err := collection.UpdateMany(ctx, condition, update)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating season in the Mongo database")
	}
	return updatedSeason, nil
}

func (m *MongoDatabase) UpdateShortCelebritiesInShow(ctx context.Context, updatedCelebrity *dto.ShortCelebrityDTO, celebrityType string) (*dto.ShortCelebrityDTO, error) {
	collection := m.client.Database("Project").Collection("Shows")
	condition := bson.D{bson.E{Key: celebrityType + ".id", Value: updatedCelebrity.ID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: celebrityType + ".$.name", Value: updatedCelebrity.Name}}},
		{Key: "$push", Value: bson.D{
			{Key: celebrityType + ".$.postersPath", Value: bson.D{
				{Key: "$each", Value: updatedCelebrity.PostersPath}}}}},
	}
	_, err := collection.UpdateMany(ctx, condition, update)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating celebrity in the Mongo database")
	}
	return updatedCelebrity, nil
}

func (m *MongoDatabase) ListShows(ctx context.Context) (dto.ShowsDTO, error) {
	collection := m.client.Database("Project").Collection("Shows")
	shows := dto.ShowsDTO{}

	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding all shows from the Mongo database")
	}

	for cursor.Next(ctx) {
		show := dto.ShowDTO{}
		err = cursor.Decode(&show)
		if err != nil {
			return nil, errors.Wrap(err, "Error while decoding show")
		}
		shows = append(shows, &show)
	}
	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error with the cursor")
	}
	cursor.Close(ctx)

	return shows, nil
}

func (m *MongoDatabase) UploadSeriesPosters(ctx context.Context, ID string, postersPath []string) (*dto.ShowDTO, error) {
	collection := m.client.Database("Project").Collection("Shows")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	update := bson.M{"$push": bson.M{
		"postersPath": bson.M{
			"$each": postersPath}}}

	updatedSeries := dto.ShowDTO{}

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	resp := collection.FindOneAndUpdate(ctx, filter, update, &opt)
	err := resp.Decode(&updatedSeries)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating series in the Mongo database")
	}

	return &updatedSeries, nil
}

func (m *MongoDatabase) DeleteSeriesPoster(ctx context.Context, ID string, image string) error {
	collection := m.client.Database("Project").Collection("Shows")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	posterPath := "/series/" + ID + "/" + image
	update := bson.M{"$pull": bson.M{"postersPath": posterPath}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "Error while updating deleted series poster in the Mongo database")
	}

	return nil
}

func (m *MongoDatabase) UploadMoviePosters(ctx context.Context, ID string, postersPath []string) (*dto.ShowDTO, error) {
	collection := m.client.Database("Project").Collection("Shows")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	update := bson.M{"$push": bson.M{
		"postersPath": bson.M{
			"$each": postersPath}}}

	updatedMovie := dto.ShowDTO{}

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	resp := collection.FindOneAndUpdate(ctx, filter, update, &opt)
	err := resp.Decode(&updatedMovie)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating movie in the Mongo database")
	}

	return &updatedMovie, nil
}

func (m *MongoDatabase) DeleteMoviePoster(ctx context.Context, ID string, image string) error {
	collection := m.client.Database("Project").Collection("Shows")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	posterPath := "/movie/" + ID + "/" + image
	update := bson.M{"$pull": bson.M{"postersPath": posterPath}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "Error while updating deleted movie poster in the Mongo database")
	}

	return nil
}

func (m *MongoDatabase) DeleteShortCelebritiesPostersInShow(ctx context.Context, celebrityID string, image string, celebrityType string) error {
	collection := m.client.Database("Project").Collection("Shows")
	filter := bson.D{bson.E{Key: celebrityType + ".id", Value: celebrityID}}
	posterPath := "/celebrities/" + celebrityID + "/" + image
	update := bson.M{"$pull": bson.M{celebrityType + ".$.postersPath": posterPath}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "Error while deleting celebrity poster in short celebrity in the Mongo database")
	}
	return nil
}

func (m *MongoDatabase) DeleteShortSeasonPostersInShow(ctx context.Context, seriesID string, seasonID string, image string) error {
	collection := m.client.Database("Project").Collection("Shows")
	filter := bson.D{bson.E{Key: "seasons.id", Value: seasonID}}
	posterPath := "/series/" + seriesID + "/" + seasonID + "/" + image
	update := bson.M{"$pull": bson.M{"seasons.$.postersPath": posterPath}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "Error while deleting season poster in the Mongo database")
	}
	return nil
}
