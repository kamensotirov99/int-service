package repository

import (
	"context"
	"int-service/dto"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleRepository interface {
	CreateArticle(ctx context.Context, newArticle *dto.ArticleDTO) (*dto.ArticleDTO, error)
	GetArticle(ctx context.Context, ID string) (*dto.ArticleDTO, error)
	UpdateArticle(ctx context.Context, updatedArticle *dto.ArticleDTO) (*dto.ArticleDTO, error)
	ListArticles(ctx context.Context, elementCount int) (dto.ArticlesDTO, error)
	ListArticlesByJournalist(ctx context.Context, ID string) (dto.ArticlesDTO, error)
	UploadArticlePosters(ctx context.Context, ID string, postersPath []string) (*dto.ArticleDTO, error)
	DeleteArticlePoster(ctx context.Context, ID string, image string) error
}

func (m *MongoDatabase) CreateArticle(ctx context.Context, newArticle *dto.ArticleDTO) (*dto.ArticleDTO, error) {
	collection := m.client.Database("Project").Collection("Articles")
	newArticle.PostersPath = []string{}
	_, err := collection.InsertOne(ctx, newArticle)
	if err != nil {
		return nil, errors.Wrap(err, "Error while inserting the new article in the Mongo database")
	}
	return newArticle, nil
}

func (m *MongoDatabase) GetArticle(ctx context.Context, ID string) (*dto.ArticleDTO, error) {
	collection := m.client.Database("Project").Collection("Articles")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	article := dto.ArticleDTO{}

	err := collection.FindOne(ctx, filter).Decode(&article)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding article by id from the Mongo database")
	}
	return &article, nil
}

func (m *MongoDatabase) UpdateArticle(ctx context.Context, updatedArticle *dto.ArticleDTO) (*dto.ArticleDTO, error) {
	collection := m.client.Database("Project").Collection("Articles")
	filter := bson.D{bson.E{Key: "id", Value: updatedArticle.ID}}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "title", Value: updatedArticle.Title},
			bson.E{Key: "releaseDate", Value: updatedArticle.ReleaseDate},
			bson.E{Key: "postersPath", Value: updatedArticle.PostersPath},
			bson.E{Key: "description", Value: updatedArticle.Description},
			bson.E{Key: "journalist", Value: updatedArticle.Journalist},
		}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating article in the Mongo database")
	}
	return updatedArticle, nil
}

func (m *MongoDatabase) ListArticles(ctx context.Context, elementCount int) (dto.ArticlesDTO, error) {
	collection := m.client.Database("Project").Collection("Articles")
	articles := dto.ArticlesDTO{}
	opts := options.FindOptions{}

	if elementCount == 0 {
		opts = *options.Find().SetSort(bson.D{{Key:"releaseDate",Value: -1}})
	} else {
		opts = *options.Find().SetSort(bson.D{{Key:"releaseDate",Value: -1}}).SetLimit(int64(elementCount))
	}
	
	cursor, err := collection.Find(ctx, bson.D{{}}, &opts)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding all articles from the Mongo database")
	}
	for cursor.Next(ctx) {
		article := dto.ArticleDTO{}
		err = cursor.Decode(&article)
		if err != nil {
			return nil, errors.Wrap(err, "Error while decoding article")
		}
		articles = append(articles, &article)

	}
	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error with the cursor")
	}
	cursor.Close(ctx)
	return articles, nil
}

func (m *MongoDatabase) ListArticlesByJournalist(ctx context.Context, journalistID string) (dto.ArticlesDTO, error) {
	collection := m.client.Database("Project").Collection("Articles")
	filter := bson.D{bson.E{Key: "journalist.id", Value: journalistID}}
	articles := dto.ArticlesDTO{}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding all articles by journalist id from the Mongo database")
	}
	for cursor.Next(ctx) {
		article := dto.ArticleDTO{}
		err = cursor.Decode(&article)
		if err != nil {
			return nil, errors.Wrap(err, "Error while decoding article")
		}
		articles = append(articles, &article)
	}
	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error with the cursor")
	}
	cursor.Close(ctx)

	return articles, nil

}

func (m *MongoDatabase) UploadArticlePosters(ctx context.Context, ID string, postersPath []string) (*dto.ArticleDTO, error) {
	collection := m.client.Database("Project").Collection("Articles")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	update := bson.M{"$push": bson.M{
		"postersPath": bson.M{
			"$each": postersPath}}}

	updatedArticle := dto.ArticleDTO{}

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	resp := collection.FindOneAndUpdate(ctx, filter, update, &opt)
	err := resp.Decode(&updatedArticle)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating article in the Mongo database")
	}

	return &updatedArticle, nil
}

func (m *MongoDatabase) DeleteArticlePoster(ctx context.Context, ID string, image string) error {
	collection := m.client.Database("Project").Collection("Articles")
	filter := bson.D{bson.E{Key: "id", Value: ID}}
	posterPath := "/articles/" + ID + "/" + image
	update := bson.M{"$pull": bson.M{"postersPath": posterPath}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "Error while updating deleted article poster in the Mongo database")
	}

	return nil
}
