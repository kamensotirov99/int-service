package repository

import (
	"context"
	"int-service/dto"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EpisodeRepository interface {
	CreateEpisode(ctx context.Context, newEpisode *dto.EpisodeDTO) (*dto.EpisodeDTO, error)
	GetEpisode(ctx context.Context, ID string) (*dto.EpisodeDTO, error)
	UpdateEpisode(ctx context.Context, updatedEpisode *dto.EpisodeDTO) (*dto.EpisodeDTO, error)
	UpdateShortCelebritiesInEpisode(ctx context.Context, updatedCelebrity *dto.ShortCelebrityDTO, celebrityType string) (*dto.ShortCelebrityDTO, error)
	UploadEpisodePosters(ctx context.Context, episodeID string, postersPath []string) (*dto.EpisodeDTO, error)
	DeleteEpisodePoster(ctx context.Context, seriesID string, seasonID string, episodeID string, image string) error
	DeleteShortCelebritiesPostersInEpisode(ctx context.Context, celebrityID string, image string, celebrityType string) error
	ListSeasonEpisodes(ctx context.Context, seasonID string) (dto.EpisodesDTO, error)
	ListCollectionEpisodes(ctx context.Context) (dto.EpisodesDTO, error)
}

func (m *MongoDatabase) CreateEpisode(ctx context.Context, newEpisode *dto.EpisodeDTO) (*dto.EpisodeDTO, error) {
	collection := m.client.Database("Project").Collection("Episodes")
	newEpisode.PostersPath = []string{}
	_, err := collection.InsertOne(ctx, newEpisode)
	if err != nil {
		return nil, errors.Wrap(err, "Error while new episode in the Mongo database")
	}
	return newEpisode, nil
}

func (m *MongoDatabase) GetEpisode(ctx context.Context, ID string) (*dto.EpisodeDTO, error) {
	collection := m.client.Database("Project").Collection("Episodes")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	episode := dto.EpisodeDTO{}

	err := collection.FindOne(ctx, filter).Decode(&episode)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding episode by id from the Mongo database")
	}
	return &episode, nil
}

func (m *MongoDatabase) UpdateEpisode(ctx context.Context, updatedEpisode *dto.EpisodeDTO) (*dto.EpisodeDTO, error) {
	collection := m.client.Database("Project").Collection("Episodes")
	filter := bson.D{bson.E{Key: "id", Value: updatedEpisode.ID}}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "title", Value: updatedEpisode.Title},
			bson.E{Key: "trailerUrl", Value: updatedEpisode.TrailerURL},
			bson.E{Key: "postersPath", Value: updatedEpisode.PostersPath},
			bson.E{Key: "length", Value: updatedEpisode.Length},
			bson.E{Key: "rating", Value: updatedEpisode.Rating},
			bson.E{Key: "resume", Value: updatedEpisode.Resume},
			bson.E{Key: "writtenBy", Value: updatedEpisode.WrittenBy},
			bson.E{Key: "producedBy", Value: updatedEpisode.ProducedBy},
			bson.E{Key: "directedBy", Value: updatedEpisode.DirectedBy},
			bson.E{Key: "starring", Value: updatedEpisode.Starring},
		}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating episode in the Mongo database")
	}

	return updatedEpisode, nil
}

func (m *MongoDatabase) UpdateShortCelebritiesInEpisode(ctx context.Context, updatedCelebrity *dto.ShortCelebrityDTO, celebrityType string) (*dto.ShortCelebrityDTO, error) {
	collection := m.client.Database("Project").Collection("Episodes")
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

func (m *MongoDatabase) UploadEpisodePosters(ctx context.Context, episodeID string, postersPath []string) (*dto.EpisodeDTO, error) {
	collection := m.client.Database("Project").Collection("Episodes")
	filter := bson.D{bson.E{Key: "id", Value: episodeID}}
	update := bson.M{"$push": bson.M{
		"postersPath": bson.M{
			"$each": postersPath}}}

	updatedEpisode := dto.EpisodeDTO{}

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	resp := collection.FindOneAndUpdate(ctx, filter, update, &opt)
	err := resp.Decode(&updatedEpisode)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating episode posters in the Mongo database")
	}

	return &updatedEpisode, nil
}

func (m *MongoDatabase) DeleteEpisodePoster(ctx context.Context, seriesID string, seasonID string, episodeID string, image string) error {
	collection := m.client.Database("Project").Collection("Episodes")
	filter := bson.D{bson.E{Key: "id", Value: episodeID}}
	posterPath := "/series/" + seriesID + "/" + seasonID + "/" + episodeID + "/" + image
	update := bson.M{"$pull": bson.M{"postersPath": posterPath}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "Error while updating deleted episode poster in the Mongo database")
	}

	return nil
}

func (m *MongoDatabase) DeleteShortCelebritiesPostersInEpisode(ctx context.Context, celebrityID string, image string, celebrityType string) error {
	collection := m.client.Database("Project").Collection("Episodes")
	filter := bson.D{bson.E{Key: celebrityType + ".id", Value: celebrityID}}
	posterPath := "/celebrities/" + celebrityID + "/" + image
	update := bson.M{"$pull": bson.M{celebrityType + ".$.postersPath": posterPath}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "Error while deleting celebrity poster in short celebrity in the Mongo database")
	}

	return nil
}

func (m *MongoDatabase) ListSeasonEpisodes(ctx context.Context, seasonID string) (dto.EpisodesDTO, error) {
	collection := m.client.Database("Project").Collection("Episodes")
	filter := bson.D{bson.E{Key: "seasonId", Value: seasonID}}
	episodes := dto.EpisodesDTO{}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding all episodes from the Mongo database")
	}
	for cursor.Next(ctx) {
		episode := dto.EpisodeDTO{}
		err = cursor.Decode(&episode)
		if err != nil {
			return nil, errors.Wrap(err, "Error while decoding episode")
		}
		episodes = append(episodes, &episode)
	}
	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error with the cursor")
	}
	cursor.Close(ctx)

	return episodes, nil
}

func (m *MongoDatabase) ListCollectionEpisodes(ctx context.Context) (dto.EpisodesDTO, error) {
	collection := m.client.Database("Project").Collection("Episodes")
	episodes := dto.EpisodesDTO{}

	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding all episodes from the Mongo database")
	}

	for cursor.Next(ctx) {
		episode := dto.EpisodeDTO{}
		err = cursor.Decode(&episode)
		if err != nil {
			return nil, errors.Wrap(err, "Error while decoding episode")
		}
		episodes = append(episodes, &episode)
	}
	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error with the cursor")
	}
	cursor.Close(ctx)

	return episodes, nil
}
