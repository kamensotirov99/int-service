package dto

import (
	"encoding/xml"
	"int-service/models"
	"time"
)

type ClothesDTO struct {
	XMLName xml.Name      `json:"-" xml:"ClothesDTO" bson:"-"`
	Clothes []ClothingDTO `json:"clothes" xml:"ClothingDTO" bson:"clothes"`
}

type ClothingDTO struct {
	XMLName xml.Name `json:"-" xml:"ClothingDTO" bson:"-"`
	ID      string   `json:"id" xml:"id" bson:"id"`
	Type    string   `json:"type" xml:"type" bson:"type"`
	Price   int      `json:"price" xml:"price" bson:"price"`
	Size    int      `json:"size" xml:"size" bson:"size"`
	Gender  string   `json:"gender" xml:"gender" bson:"gender"`
}

func (c *ClothingDTO) ToModel() *models.Clothing {
	return &models.Clothing{
		ID:     c.ID,
		Type:   c.Type,
		Price:  c.Price,
		Size:   c.Size,
		Gender: c.Gender,
	}
}

func (c *ClothesDTO) ToModel() []*models.Clothing {
	clothes := []*models.Clothing{}
	for _, c := range c.Clothes {
		clothes = append(clothes, &models.Clothing{
			ID:     c.ID,
			Type:   c.Type,
			Price:  c.Price,
			Size:   c.Size,
			Gender: c.Gender,
		})
	}
	return clothes
}

/*const (
	male   = "Male"
	female = "Female"
)*/

type GenderDTO string

type ShowsDTO []*ShowDTO

type ShowDTO struct {
	ID          string              `bson:"id"`
	Title       string              `bson:"title"`
	Type        string              `bson:"type"`
	PostersPath []string            `bson:"postersPath"`
	ReleaseDate time.Time           `bson:"releaseDate"`
	EndDate     time.Time           `bson:"endDate"`
	Rating      float64             `bson:"rating"`
	Length      ShowLengthDTO       `bson:"length"`
	TrailerURL  string              `bson:"trailerUrl"`
	Genres      ShortGenresDTO      `bson:"genres"`
	DirectedBy  FilmCrewsDTO        `bson:"directedBy"`
	ProducedBy  FilmCrewsDTO        `bson:"producedBy"`
	WrittenBy   FilmCrewsDTO        `bson:"writtenBy"`
	Starring    ShortCelebritiesDTO `bson:"starring"`
	Description string              `bson:"description"`
	Seasons     ShortSeasonsDTO     `bson:"seasons"`
}

type ShortGenresDTO []*ShortGenreDTO

type ShortGenreDTO struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

type ShortSeasonsDTO []*ShortSeasonDTO

type ShortSeasonDTO struct {
	ID          string   `bson:"id"`
	Title       string   `bson:"title"`
	PostersPath []string `bson:"postersPath"`
	Rating      float64  `bson:"rating"`
}

type SeasonsDTO []*SeasonDTO

type SeasonDTO struct {
	ID          string           `bson:"id"`
	ShowID      string           `bson:"showId"`
	Title       string           `bson:"title"`
	TrailerURL  string           `bson:"trailerUrl"`
	PostersPath []string         `bson:"postersPath"`
	Resume      string 	         `bson:"resume"`
	Rating      float64          `bson:"rating"`
	ReleaseDate time.Time        `bson:"releaseDate"`
	WrittenBy   FilmCrewsDTO     `bson:"writtenBy"`
	ProducedBy  FilmCrewsDTO     `bson:"producedBy"`
	DirectedBy  FilmCrewsDTO     `bson:"directedBy"`
	Episodes    ShortEpisodesDTO `bson:"episodes"`
}

type ShortEpisodesDTO []*ShortEpisodeDTO

type ShortEpisodeDTO struct {
	ID          string   `bson:"id"`
	Title       string   `bson:"title"`
	PostersPath []string `bson:"postersPath"`
	Rating      float64  `bson:"rating"`
	Resume      string   `bson:"resume"`
}

type EpisodesDTO []*EpisodeDTO

type EpisodeDTO struct {
	ID          string              `bson:"id"`
	SeasonID    string              `bson:"seasonId"`
	Title       string              `bson:"title"`
	PostersPath []string            `bson:"postersPath"`
	TrailerURL  string              `bson:"trailerUrl"`
	Length      ShowLengthDTO       `bson:"length"`
	Rating      float64             `bson:"rating"`
	Resume      string              `bson:"resume"`
	WrittenBy   FilmCrewsDTO        `bson:"writtenBy"`
	ProducedBy  FilmCrewsDTO        `bson:"producedBy"`
	DirectedBy  FilmCrewsDTO        `bson:"directedBy"`
	Starring    ShortCelebritiesDTO `bson:"starring"`
}

type ShortCelebritiesDTO []*ShortCelebrityDTO

type ShortCelebrityDTO struct {
	ID          string   `bson:"id"`
	Name        string   `bson:"name"`
	RoleName    string   `bson:"roleName"`
	PostersPath []string `bson:"postersPath"`
}

type FilmCrewsDTO []*FilmCrewDTO

type FilmCrewDTO struct {
	ID          string   `bson:"id"`
	Name        string   `bson:"name"`
	PostersPath []string `bson:"postersPath"`
}

type CelebritiesDTO []*CelebrityDTO

type CelebrityDTO struct {
	ID           string    `bson:"id"`
	Name         string    `bson:"name"`
	Occupation   []string  `bson:"occupation"`
	PostersPath  []string  `bson:"postersPath"`
	DateOfBirth  time.Time `bson:"dateOfBirth"`
	DateOfDeath  time.Time `bson:"dateOfDeath"`
	PlaceOfBirth string    `bson:"placeOfBirth"`
	Gender       GenderDTO `bson:"gender"`
	Bio          string    `bson:"bio"`
}

type ArticlesDTO []*ArticleDTO

type ArticleDTO struct {
	ID          string             `bson:"id"`
	Title       string             `bson:"title"`
	ReleaseDate time.Time          `bson:"releaseDate"`
	PostersPath []string           `bson:"postersPath"`
	Description string             `bson:"description"`
	Journalist  ShortJournalistDTO `bson:"journalist"`
}

type JournalistsDTO []*JournalistDTO

type JournalistDTO struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

type ShortJournalistDTO struct {
	ID string `bson:"id"`
}

type GenresDTO []*GenreDTO

type GenreDTO struct {
	ID          string `bson:"id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`
}

type ShowLengthDTO struct {
	Hours   int `bson:"hours"`
	Minutes int `bson:"minutes"`
}
