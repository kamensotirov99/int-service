package models

import (
	pb "int-service/_proto"
	"time"
)

type Clothing struct {
	ID     string
	Type   string
	Price  int
	Size   int
	Gender string
}

func (c *Clothing) ToGrpc() interface{} {
	return &pb.Clothing{
		Id:     c.ID,
		Type:   c.Type,
		Size:   int32(c.Size),
		Price:  int32(c.Price),
		Gender: c.Gender,
	}
}

type Gender string

type Shows []*Show

type Show struct {
	ID          string
	Title       string
	Type        string
	PostersPath []string
	ReleaseDate time.Time
	EndDate     time.Time
	Rating      float64
	Length      ShowLength
	TrailerURL  string
	Genres      ShortGenres
	DirectedBy  FilmCrews
	ProducedBy  FilmCrews
	WrittenBy   FilmCrews
	Starring    ShortCelebrities
	Description string
	Seasons     ShortSeasons
}

type ShortGenres []*ShortGenre

type ShortGenre struct {
	ID   string
	Name string
}

type ShortSeasons []*ShortSeason

type ShortSeason struct {
	ID          string
	Title       string
	PostersPath []string
	Rating      float64
}

type ShortCelebrities []*ShortCelebrity

type ShortCelebrity struct {
	ID          string
	Name        string
	RoleName    string
	PostersPath []string
}

type FilmCrews []*FilmCrew

type FilmCrew struct {
	ID          string
	Name        string
	PostersPath []string
}

type Genres []*Genre

type Genre struct {
	ID          string
	Name        string
	Description string
}

type ShowLength struct {
	Hours   int
	Minutes int
}

type Season struct {
	ID          string
	ShowID      string
	Title       string
	TrailerURL  string
	PostersPath []string
	Resume      string
	Rating      float64
	ReleaseDate time.Time
	WrittenBy   FilmCrews
	ProducedBy  FilmCrews
	DirectedBy  FilmCrews
	Episodes    ShortEpisodes
}

type ShortEpisodes []*ShortEpisode

type ShortEpisode struct {
	ID          string
	Title       string
	PostersPath []string
	Rating      float64
	Resume      string
}

type Episode struct {
	ID          string
	SeasonID    string
	Title       string
	PostersPath []string
	TrailerURL  string
	Length      ShowLength
	Rating      float64
	Resume      string
	WrittenBy   FilmCrews
	ProducedBy  FilmCrews
	DirectedBy  FilmCrews
	Starring    ShortCelebrities
}

type Celebrities []*Celebrity

type Celebrity struct {
	ID           string
	Name         string
	Occupation   []string
	PostersPath  []string
	DateOfBirth  time.Time
	DateOfDeath  time.Time
	PlaceOfBirth string
	Gender       Gender
	Bio          string
}

type Articles []*Article

type Article struct {
	ID          string
	Title       string
	ReleaseDate time.Time
	PostersPath []string
	Description string
	Journalist  ShortJournalist
}

type Journalists []*Journalist

type Journalist struct {
	ID   string
	Name string
}

type ShortJournalist struct {
	ID string
}
