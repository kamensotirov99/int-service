package grpc

import (
	"context"
	"errors"
	pb "int-service/_proto"
	"int-service/models"
	"time"
)

func (s *GrpcServerProject) CreateCelebrity(ctx context.Context, req *pb.CreateCelebrityRequest) (*pb.Celebrity, error) {
	dateOfBirth := req.DateOfBirth.AsTime()
	if err := req.DateOfBirth.CheckValid(); err != nil {
		return nil, errors.New("error while converting dateOfBirth")
	}
	var dateOfDeath time.Time
	if req.DateOfDeath != nil {
		if err := req.DateOfDeath.CheckValid(); err != nil {
			return nil, errors.New("error while converting dateOfDeath")
		}
	}

	gender := models.Gender(req.Gender)
	resp, err := s.service.CreateCelebrity(ctx, req.Name, req.Occupation, req.PostersPath, dateOfBirth, dateOfDeath, req.PlaceOfBirth, &gender, req.Bio)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Celebrity), nil
}

func (s *GrpcServerProject) GetCelebrity(ctx context.Context, req *pb.GetByIDRequest) (*pb.Celebrity, error) {
	resp, err := s.service.GetCelebrity(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return resp.ToGrpc().(*pb.Celebrity), nil
}

func (s *GrpcServerProject) UpdateCelebrity(ctx context.Context, req *pb.Celebrity) (*pb.Celebrity, error) {
	dateOfBirth := req.DateOfBirth.AsTime()
	if err := req.DateOfBirth.CheckValid(); err != nil {
		return nil, errors.New("error while converting dateOfBirth")
	}
	var dateOfDeath time.Time
	if req.DateOfDeath != nil {
		if err := req.DateOfDeath.CheckValid(); err != nil {
			return nil, errors.New("error while converting dateOfDeath")
		}
	}

	gender := models.Gender(req.Gender)

	resp, err := s.service.UpdateCelebrity(ctx, req.Id, req.Name, req.Occupation, req.PostersPath, dateOfBirth, dateOfDeath, req.PlaceOfBirth, &gender, req.Bio)
	if err != nil {
		return nil, errors.New("error while updating celebrity")
	}
	return resp.ToGrpc().(*pb.Celebrity), nil
}

func (s *GrpcServerProject) UploadCelebrityPosters(ctx context.Context, req *pb.UploadCelebrityPostersRequest) (*pb.Celebrity, error) {
	resp, err := s.service.UploadCelebrityPosters(ctx, req.CelebrityId, req.PostersPath)
	if err != nil {
		return nil, errors.New("error while uploading celebrity posters")
	}
	return resp.ToGrpc().(*pb.Celebrity), nil
}

func (s *GrpcServerProject) DeleteCelebrityPoster(ctx context.Context, req *pb.DeleteCelebrityPosterRequest) (*pb.EmptyResponse, error) {
	err := s.service.DeleteCelebrityPoster(ctx, req.CelebrityId, req.Image)
	if err != nil {
		return nil, errors.New("error while deleting celebrity poster")
	}
	return &pb.EmptyResponse{}, nil
}

func (s *GrpcServerProject) ListCelebrities(ctx context.Context, req *pb.GetAllRequest) (*pb.CelebrityListResponse, error) {
	resp, err := s.service.ListCelebrities(ctx)
	if err != nil {
		return nil, err
	}

	celebs := &pb.CelebrityListResponse{}
	for _, celeb := range resp {
		celebs.Celebrities = append(celebs.Celebrities, celeb.ToGrpc().(*pb.Celebrity))
	}
	return celebs, nil
}

func toFilmCrewsModel(filmCrewsPb *pb.FilmCrew) models.FilmCrews {
	filmCrews := models.FilmCrews{}
	for _, filmCrew := range filmCrewsPb.FilmCrew {
		filmCrews = append(filmCrews, &models.FilmCrew{
			ID:          filmCrew.Id,
			Name:        filmCrew.Name,
			PostersPath: filmCrew.PostersPath,
		})
	}
	return filmCrews
}

func toShortCelebsModel(shortCelebsPb *pb.ShortCelebrities) models.ShortCelebrities {
	shortCelebs := models.ShortCelebrities{}
	for _, shortCeleb := range shortCelebsPb.ShortCelebs {
		shortCelebs = append(shortCelebs, &models.ShortCelebrity{
			ID:          shortCeleb.Id,
			Name:        shortCeleb.Name,
			RoleName:    shortCeleb.RoleName,
			PostersPath: shortCeleb.PostersPath,
		})
	}
	return shortCelebs
}
