package service

import (
	"context"
	"int-service/dto"
	"int-service/models"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	Actor      = "Actor"
	Actress    = "Actress"
	Writer     = "Writer"
	Director   = "Director"
	Producer   = "Producer"
	starring   = "starring"
	writtenBy  = "writtenBy"
	directedBy = "directedBy"
	producedBy = "producedBy"
)

type CelebrityServicer interface {
	CreateCelebrity(ctx context.Context, name string, occupation []string, postersPath []string, dateOfBirth time.Time, dateOfDeath time.Time, placeOfBirth string, genderModel *models.Gender, bio string) (models.ResponseModeler, error)
	GetCelebrity(ctx context.Context, ID string) (models.ResponseModeler, error)
	UpdateCelebrity(ctx context.Context, ID string, name string, occupation []string, postersPath []string, dateOfBirth time.Time, dateOfDeath time.Time, placeOfBirth string, genderModel *models.Gender, bio string) (models.ResponseModeler, error)
	UploadCelebrityPosters(ctx context.Context, ID string, postersPath []string) (models.ResponseModeler, error)
	DeleteCelebrityPoster(ctx context.Context, ID string, image string) error
	ListCelebrities(ctx context.Context) ([]models.ResponseModeler, error)
}

func (s *projectService) CreateCelebrity(ctx context.Context, name string, occupation []string, postersPath []string, dateOfBirth time.Time, dateOfDeath time.Time, placeOfBirth string, genderModel *models.Gender, bio string) (models.ResponseModeler, error) {
	err := s.validateCelebrityUniqueness(ctx, name, dateOfBirth)
	if err != nil {
		s.logger.Error("Error while creating celebrity")
		return nil, errors.Wrap(err, "Error while creating celebrity")
	}
	celebrity := toCelebrityDTO(uuid.New().String(), name, occupation, postersPath, dateOfBirth, dateOfDeath, placeOfBirth, genderModel, bio)
	resp, err := s.repository.CreateCelebrity(ctx, celebrity)
	if err != nil {
		s.logger.Error("Error while creating celebrity")
		return nil, errors.Wrap(err, "Error while creating celebrity")
	}
	return resp.ToModel(), nil
}

func (s *projectService) GetCelebrity(ctx context.Context, ID string) (models.ResponseModeler, error) {
	resp, err := s.repository.GetCelebrity(ctx, ID)
	if err != nil {
		s.logger.Error("Error while getting celebrity by id")
		return nil, errors.Wrap(err, "Error while getting celebrity by id")
	}
	return resp.ToModel(), nil
}

func (s *projectService) UpdateCelebrity(ctx context.Context, ID string, name string, occupation []string, postersPath []string, dateOfBirth time.Time, dateOfDeath time.Time, placeOfBirth string, genderModel *models.Gender, bio string) (models.ResponseModeler, error) {
	updatedCelebrity := toCelebrityDTO(ID, name, occupation, postersPath, dateOfBirth, dateOfDeath, placeOfBirth, genderModel, bio)
	resp, err := s.repository.UpdateCelebrity(ctx, updatedCelebrity)
	if err != nil {
		s.logger.Error("Error while updating celebrity")
		return nil, errors.New("Error while updating celebrity")
	}
	shortCeleb := &dto.ShortCelebrityDTO{
		ID:   resp.ID,
		Name: resp.Name,
	}
	err = s.updateShortCelebrities(ctx, shortCeleb, occupation)
	if err != nil {
		s.logger.Error("Error while updating short celebrity")
		return nil, errors.Wrap(err, "Error while updating short celebrity")
	}
	return resp.ToModel(), nil
}

func (s *projectService) UploadCelebrityPosters(ctx context.Context, ID string, postersPath []string) (models.ResponseModeler, error) {
	resp, err := s.repository.UploadCelebrityPosters(ctx, ID, postersPath)
	if err != nil {
		s.logger.Error("Error while uploading celebrity posters")
		return nil, errors.Wrap(err, "Error while uploading celebrity posters")
	}
	celeb, err := s.repository.GetCelebrity(ctx, ID)
	if err != nil {
		s.logger.Error("Error while getting celebrity by id")
		return nil, errors.Wrap(err, "Error while getting celebrity by id")
	}
	shortCeleb := &dto.ShortCelebrityDTO{
		ID:          ID,
		Name:        celeb.Name,
		PostersPath: postersPath,
	}
	err = s.updateShortCelebrities(ctx, shortCeleb, celeb.Occupation)
	if err != nil {
		s.logger.Error("Error while updating short celebrity posters")
		return nil, errors.Wrap(err, "Error while updating short celebrity posters")
	}
	return resp.ToModel(), nil
}

func (s *projectService) DeleteCelebrityPoster(ctx context.Context, ID string, image string) error {
	err := s.repository.DeleteCelebrityPoster(ctx, ID, image)
	if err != nil {
		s.logger.Error("Error while deleting celebrity poster in database")
		return errors.Wrap(err, "Error while deleting celebrity poster in database")
	}
	celeb, err := s.repository.GetCelebrity(ctx, ID)
	if err != nil {
		s.logger.Error("Error while getting celebrity by id")
		errors.Wrap(err, "Error while getting celebrity by id")
	}
	err = s.deleteShortCelebritiesPosters(ctx, ID, image, celeb.Occupation)
	if err != nil {
		s.logger.Error("Error while deleting short celebrity posters")
		errors.Wrap(err, "Error while deleting short celebrity posters")
	}
	return nil
}

func (s *projectService) ListCelebrities(ctx context.Context) ([]models.ResponseModeler, error) {
	resp, err := s.repository.ListCelebrities(ctx)
	if err != nil {
		s.logger.Error("Error while listing all celebrities")
		return nil, errors.New("Error while listing all celebrities")
	}
	celebrities := []models.ResponseModeler{}
	for _, celebrity := range resp {
		celebrities = append(celebrities, celebrity.ToModel())
	}
	return celebrities, nil
}

func (s *projectService) validateCelebrityUniqueness(ctx context.Context, name string, dateOfBirth time.Time) error {
	celebs, err := s.repository.ListCelebrities(ctx)
	if err != nil {
		return errors.Wrap(err, "Error while getting celebrities from the Mongo database")
	}
	for _, celeb := range celebs {
		if celeb.Name == name && celeb.DateOfBirth == dateOfBirth {
			return errors.New("There is already a celebrity with that name and date of birth.")
		}
	}
	return nil
}

func (s *projectService) updateShortCelebrities(ctx context.Context, shortCeleb *dto.ShortCelebrityDTO, occupation []string) error {
	celebrityType := ""
	for _, occupationType := range occupation {
		if occupationType == Actress || occupationType == Actor {
			celebrityType = starring
		} else if occupationType == Writer {
			celebrityType = writtenBy
		} else if occupationType == Director {
			celebrityType = directedBy
		} else if occupationType == Producer {
			celebrityType = producedBy
		}
		if _, err := s.repository.UpdateShortCelebritiesInShow(ctx, shortCeleb, celebrityType); err != nil {
			s.logger.Error("Error while updating short " + celebrityType + " in show")
			return errors.New("Error while updating short " + celebrityType + " in show")
		}
		if _, err := s.repository.UpdateShortCelebritiesInEpisode(ctx, shortCeleb, celebrityType); err != nil {
			s.logger.Error("Error while updating short " + celebrityType + " in episode")
			return errors.New("Error while updating short " + celebrityType + " in episode")
		}
		if _, err := s.repository.UpdateShortCelebritiesInSeasons(ctx, shortCeleb, celebrityType); err != nil {
			s.logger.Error("Error while updating short " + celebrityType + " in season")
			return errors.New("Error while updating short " + celebrityType + " in season")
		}
	}
	return nil
}

func (s *projectService) deleteShortCelebritiesPosters(ctx context.Context, ID string, image string, occupation []string) error {
	celebrityType := ""
	for _, occupationType := range occupation {
		if occupationType == Actress || occupationType == Actor {
			celebrityType = starring
		} else if occupationType == Writer {
			celebrityType = writtenBy
		} else if occupationType == Director {
			celebrityType = directedBy
		} else if occupationType == Producer {
			celebrityType = producedBy
		}
		if err := s.repository.DeleteShortCelebritiesPostersInShow(ctx, ID, image, celebrityType); err != nil {
			s.logger.Error("Error while updating short " + celebrityType + " in show")
			return errors.New("Error while updating short " + celebrityType + " in show")
		}
		if err := s.repository.DeleteShortCelebritiesPostersInEpisode(ctx, ID, image, celebrityType); err != nil {
			s.logger.Error("Error while updating short " + celebrityType + " in episode")
			return errors.New("Error while updating short " + celebrityType + " in episode")
		}
		if err := s.repository.DeleteShortCelebritiesPostersInSeason(ctx, ID, image, celebrityType); err != nil {
			s.logger.Error("Error while updating short " + celebrityType + " in season")
			return errors.New("Error while updating short " + celebrityType + " in season")
		}
	}
	return nil
}

func toCelebrityDTO(ID string, name string, occupation []string, postersPath []string, dateOfBirth time.Time, dateOfDeath time.Time, placeOfBirth string, genderModel *models.Gender, bio string) *dto.CelebrityDTO {
	gender := dto.GenderDTO(*genderModel)
	celebrity := dto.CelebrityDTO{
		ID:           ID,
		Name:         name,
		Occupation:   occupation,
		PostersPath:  postersPath,
		DateOfBirth:  dateOfBirth,
		DateOfDeath:  dateOfDeath,
		PlaceOfBirth: placeOfBirth,
		Gender:       gender,
		Bio:          bio,
	}
	return &celebrity
}

func toFilmCrewsDTO(filmCrewsModel models.FilmCrews) dto.FilmCrewsDTO {
	filmCrews := dto.FilmCrewsDTO{}
	for _, filmCrew := range filmCrewsModel {
		filmCrewDTO := dto.FilmCrewDTO{
			ID:          filmCrew.ID,
			Name:        filmCrew.Name,
			PostersPath: filmCrew.PostersPath,
		}
		filmCrews = append(filmCrews, &filmCrewDTO)
	}
	return filmCrews
}

func toShortCelebDTO(shortCelebModel models.ShortCelebrities) dto.ShortCelebritiesDTO {
	starring := dto.ShortCelebritiesDTO{}
	for _, celeb := range shortCelebModel {
		celebDTO := dto.ShortCelebrityDTO{
			ID:          celeb.ID,
			Name:        celeb.Name,
			RoleName:    celeb.RoleName,
			PostersPath: celeb.PostersPath,
		}
		starring = append(starring, &celebDTO)
	}
	return starring
}
