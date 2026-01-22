package models

import (
	"time"
)

type Movie struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Year        int       `json:"year"`
	Rating      float64   `json:"rating"`
	Duration    int       `json:"duration"` // in minutes
	Genre       string    `json:"genre"`
	Director    string    `json:"director"`
	PosterURL   string    `json:"poster_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Actors      []Actor   `json:"actors,omitempty"`
	Reviews     []Review  `json:"reviews,omitempty"`
}

type Actor struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	BirthDate  string `json:"birth_date"`
	Nationality string `json:"nationality"`
	Biography  string `json:"biography"`
	ProfileURL string `json:"profile_url"`
}

type Review struct {
	ID        string    `json:"id"`
	MovieID   string    `json:"movie_id"`
	UserName  string    `json:"user_name"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type MovieActor struct {
	MovieID       string `json:"movie_id"`
	ActorID       string `json:"actor_id"`
	CharacterName string `json:"character_name"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type MovieFilter struct {
	Genre   string  `json:"genre"`
	MinYear int     `json:"min_year"`
	MaxYear int     `json:"max_year"`
	MinRating float64 `json:"min_rating"`
	Search  string  `json:"search"`
}