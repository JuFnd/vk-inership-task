package communication

import "filmoteka/pkg/models"

type (
	FilmsListResponse struct {
		Films []models.FilmItem `json:"films"`
	}

	FindFilmResponse struct {
		Films []models.FilmShortItem `json:"film_data"`
	}

	ActorsListResponse struct {
		Actors []models.ActorItem `json:"actors"`
	}
)
