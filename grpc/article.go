package grpc

import (
	"context"
	"errors"
	pb "int-service/_proto"
	"int-service/models"
)

func (s *GrpcServerProject) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.Article, error) {
	releaseDate := req.ReleaseDate.AsTime()
	if err := req.ReleaseDate.CheckValid(); err != nil {
		return nil, errors.New("error while converting releaseDate")
	}
	journalist := models.Journalist{
		Name: req.Journalist.Name,
	}
	resp, err := s.service.CreateArticle(ctx, req.Title, releaseDate, req.PostersPath, req.Description, journalist.Name)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Article), nil
}

func (s *GrpcServerProject) GetArticle(ctx context.Context, req *pb.GetByIDRequest) (*pb.Article, error) {
	resp, err := s.service.GetArticle(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Article), nil
}

func (s *GrpcServerProject) UpdateArticle(ctx context.Context, req *pb.Article) (*pb.Article, error) {
	releaseDate := req.ReleaseDate.AsTime()
	if err := req.ReleaseDate.CheckValid(); err != nil {
		return nil, errors.New("error while converting releaseDate")
	}
	journalist := models.Journalist{
		ID: req.Journalist.Id,
	}
	resp, err := s.service.UpdateArticle(ctx, req.Id, req.Title, releaseDate, req.PostersPath, req.Description, &journalist)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Article), nil
}

func (s *GrpcServerProject) ListArticles(ctx context.Context, req *pb.ListArticlesRequest) (*pb.ArticleListResponse, error) {
	resp, err := s.service.ListArticles(ctx,int(req.ElementCount))
	if err != nil {
		return nil, err
	}

	articles := &pb.ArticleListResponse{}
	for _, article := range resp {
		articles.Articles = append(articles.Articles, article.ToGrpc().(*pb.Article))
	}
	return articles, nil
}

func (s *GrpcServerProject) ListArticlesByJournalist(ctx context.Context, req *pb.GetByIDRequest) (*pb.ArticleListResponse, error) {
	resp, err := s.service.ListArticlesByJournalist(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	articles := &pb.ArticleListResponse{}
	for _, article := range resp {
		articles.Articles = append(articles.Articles, article.ToGrpc().(*pb.Article))
	}
	return articles, nil
}

func (s *GrpcServerProject) UploadArticlePosters(ctx context.Context, req *pb.UploadArticlePostersRequest) (*pb.Article, error) {
	resp, err := s.service.UploadArticlePosters(ctx, req.ArticleId, req.PostersPath)
	if err != nil {
		return nil, errors.New("error while updating article posters")
	}
	return resp.ToGrpc().(*pb.Article), nil
}

func (s *GrpcServerProject) DeleteArticlePoster(ctx context.Context, req *pb.DeleteArticlePosterRequest) (*pb.EmptyResponse, error) {
	err := s.service.DeleteArticlePoster(ctx, req.ArticleId, req.Image)
	if err != nil {
		return nil, errors.New("error while updating deleted article poster")
	}
	return &pb.EmptyResponse{}, nil
}
