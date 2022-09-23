package models

import (
	pb "int-service/_proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ResponseModeler interface {
	ToGrpc() interface{}
}

func (a *Article) ToGrpc() interface{} {
	releaseDate := timestamppb.New(a.ReleaseDate)
	return &pb.Article{
		Id:          a.ID,
		Title:       a.Title,
		ReleaseDate: releaseDate,
		PostersPath: a.PostersPath,
		Description: a.Description,
		Journalist:  a.Journalist.ToGrpc().(*pb.ShortJournalist),
	}
}

func (j *Journalist) ToGrpc() interface{} {
	return &pb.Journalist{
		Id:   j.ID,
		Name: j.Name,
	}
}

func (j *ShortJournalist) ToGrpc() interface{} {
	return &pb.ShortJournalist{
		Id: j.ID,
	}
}

func (c *Celebrity) ToGrpc() interface{} {
	dateOfBirth := timestamppb.New(c.DateOfBirth)
	dateOfDeath := timestamppb.New(c.DateOfDeath)
	return &pb.Celebrity{
		Id:           c.ID,
		Name:         c.Name,
		Occupation:   c.Occupation,
		PostersPath:  c.PostersPath,
		DateOfBirth:  dateOfBirth,
		DateOfDeath:  dateOfDeath,
		PlaceOfBirth: c.PlaceOfBirth,
		Gender:       string(c.Gender),
		Bio:          c.Bio,
	}
}

func (e *Episode) ToGrpc() interface{} {
	return &pb.Episode{
		Id:          e.ID,
		SeasonId:    e.SeasonID,
		Title:       e.Title,
		PostersPath: e.PostersPath,
		TrailerUrl:  e.TrailerURL,
		ShowLength:  e.Length.ToGrpc().(*pb.ShowLength),
		Rating:      e.Rating,
		Resume:      e.Resume,
		WrittenBy:   e.WrittenBy.ToGrpc().(*pb.FilmCrew),
		ProducedBy:  e.ProducedBy.ToGrpc().(*pb.FilmCrew),
		DirectedBy:  e.DirectedBy.ToGrpc().(*pb.FilmCrew),
		Starring:    e.Starring.ToGrpc().(*pb.ShortCelebrities),
	}
}

func (s *ShowLength) ToGrpc() interface{} {
	return &pb.ShowLength{
		Hours:   int32(s.Hours),
		Minutes: int32(s.Minutes),
	}
}

func (f *FilmCrews) ToGrpc() interface{} {
	filmCrews := &pb.FilmCrew{}
	for _, filmCrew := range *f {
		filmCrews.FilmCrew = append(filmCrews.FilmCrew, filmCrew.ToGrpc().(*pb.FilmStaff))
	}
	return filmCrews
}

func (f *FilmCrew) ToGrpc() interface{} {
	return &pb.FilmStaff{
		Id:          f.ID,
		Name:        f.Name,
		PostersPath: f.PostersPath,
	}
}

func (s *ShortCelebrities) ToGrpc() interface{} {
	shortCelebs := &pb.ShortCelebrities{}
	for _, shortCeleb := range *s {
		shortCelebs.ShortCelebs = append(shortCelebs.ShortCelebs, shortCeleb.ToGrpc().(*pb.ShortCelebrity))
	}
	return shortCelebs
}

func (s *ShortCelebrity) ToGrpc() interface{} {
	return &pb.ShortCelebrity{
		Id:          s.ID,
		Name:        s.Name,
		RoleName:    s.RoleName,
		PostersPath: s.PostersPath,
	}
}

func (s *Show) ToGrpc() interface{} {
	releaseDate := timestamppb.New(s.ReleaseDate)
	endDate := timestamppb.New(s.EndDate)
	return &pb.Show{
		Id:          s.ID,
		Title:       s.Title,
		Type:        s.Type,
		PostersPath: s.PostersPath,
		ReleaseDate: releaseDate,
		EndDate:     endDate,
		Rating:      s.Rating,
		Length:      s.Length.ToGrpc().(*pb.ShowLength),
		TrailerUrl:  s.TrailerURL,
		Genres:      s.Genres.ToGrpc().(*pb.ShortGenres),
		DirectedBy:  s.DirectedBy.ToGrpc().(*pb.FilmCrew),
		ProducedBy:  s.ProducedBy.ToGrpc().(*pb.FilmCrew),
		WrittenBy:   s.WrittenBy.ToGrpc().(*pb.FilmCrew),
		Starring:    s.Starring.ToGrpc().(*pb.ShortCelebrities),
		Description: s.Description,
		Seasons:     s.Seasons.ToGrpc().(*pb.ShortSeasons),
	}
}

func (g *ShortGenres) ToGrpc() interface{} {
	genres := &pb.ShortGenres{}
	for _, genre := range *g {
		genres.Genres = append(genres.Genres, genre.ToGrpc().(*pb.ShortGenre))
	}
	return genres
}

func (g *ShortGenre) ToGrpc() interface{} {
	return &pb.ShortGenre{
		Id:   g.ID,
		Name: g.Name,
	}
}

func (g *Genres) ToGrpc() interface{} {
	genres := &pb.GenreListResponse{}
	for _, genre := range *g {
		genres.Genres = append(genres.Genres, genre.ToGrpc().(*pb.Genre))
	}
	return genres
}

func (g *Genre) ToGrpc() interface{} {
	return &pb.Genre{
		Id:          g.ID,
		Name:        g.Name,
		Description: g.Description,
	}
}

func (s *ShortSeasons) ToGrpc() interface{} {
	seasons := &pb.ShortSeasons{}
	for _, season := range *s {
		seasons.Seasons = append(seasons.Seasons, season.ToGrpc().(*pb.ShortSeason))
	}
	return seasons
}

func (s *ShortSeason) ToGrpc() interface{} {
	return &pb.ShortSeason{
		Id:          s.ID,
		Title:       s.Title,
		PostersPath: s.PostersPath,
		Rating:      s.Rating,
	}
}

func (s *Season) ToGrpc() interface{} {
	releaseDate := timestamppb.New(s.ReleaseDate)
	return &pb.Season{
		Id:          s.ID,
		ShowId:      s.ShowID,
		Title:       s.Title,
		TrailerUrl:  s.TrailerURL,
		PostersPath: s.PostersPath,
		Resume:      s.Resume,
		Rating:      s.Rating,
		ReleaseDate: releaseDate,
		DirectedBy:  s.DirectedBy.ToGrpc().(*pb.FilmCrew),
		ProducedBy:  s.ProducedBy.ToGrpc().(*pb.FilmCrew),
		WrittenBy:   s.WrittenBy.ToGrpc().(*pb.FilmCrew),
		Episodes:    s.Episodes.ToGrpc().(*pb.ShortEpisodeList),
	}
}

func (e *ShortEpisodes) ToGrpc() interface{} {
	episodes := &pb.ShortEpisodeList{}
	for _, episode := range *e {
		episodes.ShortEpisodes = append(episodes.ShortEpisodes, episode.ToGrpc().(*pb.ShortEpisode))
	}
	return episodes
}

func (e *ShortEpisode) ToGrpc() interface{} {
	return &pb.ShortEpisode{
		Id:          e.ID,
		Title:       e.Title,
		PostersPath: e.PostersPath,
		Rating:      e.Rating,
		Resume:      e.Resume,
	}
}
