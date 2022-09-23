package dto

import "int-service/models"

func (s *ShowDTO) ToModel() *models.Show {
	return &models.Show{
		ID:          s.ID,
		Title:       s.Title,
		Type:        s.Type,
		PostersPath: s.PostersPath,
		ReleaseDate: s.ReleaseDate,
		EndDate:     s.EndDate,
		Rating:      s.Rating,
		Length:      s.Length.ToModel().(models.ShowLength),
		TrailerURL:  s.TrailerURL,
		Genres:      s.Genres.ToModel().(models.ShortGenres),
		DirectedBy:  s.DirectedBy.ToModel().(models.FilmCrews),
		WrittenBy:   s.WrittenBy.ToModel().(models.FilmCrews),
		ProducedBy:  s.ProducedBy.ToModel().(models.FilmCrews),
		Description: s.Description,
		Starring:    s.Starring.ToModel().(models.ShortCelebrities),
		Seasons:     s.Seasons.ToModel().(models.ShortSeasons),
	}
}

func (g *ShortGenresDTO) ToModel() interface{} {
	genres := models.ShortGenres{}
	for _, genre := range *g {
		genres = append(genres, genre.ToModel().(*models.ShortGenre))
	}
	return genres
}

func (g *ShortGenreDTO) ToModel() interface{} {
	return &models.ShortGenre{
		ID:   g.ID,
		Name: g.Name,
	}
}

func (s *ShowLengthDTO) ToModel() interface{} {
	return models.ShowLength{
		Hours:   s.Hours,
		Minutes: s.Minutes,
	}
}

func (g *GenresDTO) ToModel() interface{} {
	genresModel := models.Genres{}
	for _, genre := range *g {
		genresModel = append(genresModel, genre.ToModel())
	}
	return genresModel
}

func (g *GenreDTO) ToModel() *models.Genre {
	return &models.Genre{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
	}
}

func (f *FilmCrewsDTO) ToModel() interface{} {
	filmCrewsModel := models.FilmCrews{}
	for _, filmCrew := range *f {
		filmCrewsModel = append(filmCrewsModel, filmCrew.ToModel().(*models.FilmCrew))
	}
	return filmCrewsModel
}

func (f *FilmCrewDTO) ToModel() interface{} {
	return &models.FilmCrew{
		ID:          f.ID,
		Name:        f.Name,
		PostersPath: f.PostersPath,
	}
}

func (s *ShortCelebritiesDTO) ToModel() interface{} {
	shortCelebsModel := models.ShortCelebrities{}
	for _, shortCeleb := range *s {
		shortCelebsModel = append(shortCelebsModel, shortCeleb.ToModel().(*models.ShortCelebrity))
	}
	return shortCelebsModel
}

func (s *ShortCelebrityDTO) ToModel() interface{} {
	return &models.ShortCelebrity{
		ID:          s.ID,
		Name:        s.Name,
		RoleName:    s.RoleName,
		PostersPath: s.PostersPath,
	}
}

func (s *ShortSeasonsDTO) ToModel() interface{} {
	shortSeasonsModel := models.ShortSeasons{}
	for _, shortSeason := range *s {
		shortSeasonsModel = append(shortSeasonsModel, shortSeason.ToModel().(*models.ShortSeason))
	}
	return shortSeasonsModel
}

func (s *ShortSeasonDTO) ToModel() interface{} {
	return &models.ShortSeason{
		ID:          s.ID,
		Title:       s.Title,
		PostersPath: s.PostersPath,
		Rating:      s.Rating,
	}
}

func (s *SeasonDTO) ToModel() *models.Season {
	return &models.Season{
		ID:          s.ID,
		ShowID:      s.ShowID,
		Title:       s.Title,
		TrailerURL:  s.TrailerURL,
		PostersPath: s.PostersPath,
		Resume:      s.Resume,
		Rating:      s.Rating,
		ReleaseDate: s.ReleaseDate,
		WrittenBy:   s.WrittenBy.ToModel().(models.FilmCrews),
		ProducedBy:  s.ProducedBy.ToModel().(models.FilmCrews),
		DirectedBy:  s.DirectedBy.ToModel().(models.FilmCrews),
		Episodes:    s.Episodes.ToModel().(models.ShortEpisodes),
	}
}

func (s *ShortEpisodesDTO) ToModel() interface{} {
	shortEpisodesModel := models.ShortEpisodes{}
	for _, shortEpisode := range *s {
		shortEpisodesModel = append(shortEpisodesModel, shortEpisode.ToModel().(*models.ShortEpisode))
	}
	return shortEpisodesModel
}

func (s *ShortEpisodeDTO) ToModel() interface{} {
	return &models.ShortEpisode{
		ID:          s.ID,
		Title:       s.Title,
		PostersPath: s.PostersPath,
		Rating:      s.Rating,
		Resume:      s.Resume,
	}
}

func (e *EpisodeDTO) ToModel() *models.Episode {
	return &models.Episode{
		ID:          e.ID,
		SeasonID:    e.SeasonID,
		Title:       e.Title,
		PostersPath: e.PostersPath,
		TrailerURL:  e.TrailerURL,
		Length:      e.Length.ToModel().(models.ShowLength),
		Rating:      e.Rating,
		Resume:      e.Resume,
		WrittenBy:   e.WrittenBy.ToModel().(models.FilmCrews),
		ProducedBy:  e.ProducedBy.ToModel().(models.FilmCrews),
		DirectedBy:  e.ProducedBy.ToModel().(models.FilmCrews),
		Starring:    e.Starring.ToModel().(models.ShortCelebrities),
	}
}

func (g *GenderDTO) ToModel() interface{} {
	return models.Gender(*g)
}

func (c *CelebrityDTO) ToModel() *models.Celebrity {
	return &models.Celebrity{
		ID:           c.ID,
		Name:         c.Name,
		Occupation:   c.Occupation,
		PostersPath:  c.PostersPath,
		DateOfBirth:  c.DateOfBirth,
		DateOfDeath:  c.DateOfDeath,
		PlaceOfBirth: c.PlaceOfBirth,
		Gender:       c.Gender.ToModel().(models.Gender),
		Bio:          c.Bio,
	}
}

func (a *ArticleDTO) ToModel() *models.Article {
	return &models.Article{
		ID:          a.ID,
		Title:       a.Title,
		ReleaseDate: a.ReleaseDate,
		PostersPath: a.PostersPath,
		Description: a.Description,
		Journalist:  *a.Journalist.ToModel(),
	}
}

func (j *JournalistDTO) ToModel() *models.Journalist {
	return &models.Journalist{
		ID:   j.ID,
		Name: j.Name,
	}
}

func (j *ShortJournalistDTO) ToModel() *models.ShortJournalist {
	return &models.ShortJournalist{
		ID: j.ID,
	}
}
