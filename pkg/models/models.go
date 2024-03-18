package models

import "time"

type (
	Session struct {
		Login     string
		SID       string
		ExpiresAt time.Time
	}

	UserItem struct {
		Login string `json:"login"`
	}

	FilmItem struct {
		Id          int         `json:"id"`
		Title       string      `json:"title"`
		Description string      `json:"description"`
		Rating      float64     `json:"rating"`
		ReleaseDate string      `json:"release_date"`
		Crew        []ActorItem `json:"crew"`
	}

	FilmShortItem struct {
		Id          int              `json:"id"`
		Title       string           `json:"title"`
		Description string           `json:"description"`
		Rating      float64          `json:"rating"`
		ReleaseDate string           `json:"release_date"`
		Crew        []ActorShortItem `json:"crew"`
	}

	ActorShortItem struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Gender    string `json:"gender"`
		BirthDate string `json:"birth_date"`
	}

	ActorItem struct {
		Id        int             `json:"id"`
		Name      string          `json:"name"`
		Gender    string          `json:"gender"`
		BirthDate string          `json:"birth_date"`
		Films     []FilmShortItem `json:"films"`
	}
)
