package delivery

import (
	"context"
	"filmoteka/pkg/middleware"
	communication "filmoteka/pkg/requests"
	"filmoteka/pkg/util"
	"filmoteka/pkg/variables"
	"log/slog"
	"net/http"
)

// Core interface

//go:generate mockgen -source=api.go -destination=../mocks/core_mock.go -package=mocks
type ICore interface {
	GetFilms(begin uint64, end uint64, sortType string) (communication.FilmsListResponse, error)
	FindFilm(filmName string, actorName string) (communication.FindFilmResponse, error)
	AddFilm(title string, description string, rating float64, releaseDate string, crew []int64) error
	EditFilm(id int64, title string, description string, rating float64, releaseDate string, crew []int64) error
	GetActors(begin uint64, end uint64) (communication.ActorsListResponse, error)
	AddActor(name string, gender string, birthdate string) error
	EditActor(id int64, name string, gender string, birthdate string, films []int64) error
	DeleteActor(id int64) error
	DeleteFilm(id int64) error
	GetUserRole(ctx context.Context, id int64) (string, error)
	GetUserId(ctx context.Context, sid string) (int64, error)
}

type API struct {
	core   ICore
	logger *slog.Logger
	mux    *http.ServeMux
}

func (api *API) ListenAndServe(appConfig *variables.AppConfig) error {
	err := http.ListenAndServe(appConfig.Address, api.mux)
	if err != nil {
		api.logger.Error(variables.ListenAndServeError, err.Error())
		return err
	}
	return nil
}

func GetFilmsApi(filmsCore ICore, filmsLogger *slog.Logger) *API {
	api := &API{
		core:   filmsCore,
		logger: filmsLogger,
		mux:    http.NewServeMux(),
	}

	// Actors handlers
	api.mux.Handle("/api/v1/actors", middleware.MethodMiddleware(
		middleware.AuthorizationMiddleware(
			http.HandlerFunc(api.GetActors),
			api.core, api.logger),
		http.MethodGet,
		api.logger))

	api.mux.Handle("/api/v1/actors/add", middleware.MethodMiddleware(
		middleware.AuthorizationMiddleware(
			middleware.PermissionsMiddleware(
				http.HandlerFunc(api.AddInfoAboutActor), api.core, variables.AdminRole, api.logger),
			api.core, api.logger),
		http.MethodPost,
		api.logger))

	api.mux.Handle("/api/v1/actors/edit", middleware.MethodMiddleware(
		middleware.AuthorizationMiddleware(
			middleware.PermissionsMiddleware(
				http.HandlerFunc(api.EditInfoAboutActor), api.core, variables.AdminRole, api.logger),
			api.core, api.logger),
		http.MethodPost,
		api.logger))

	api.mux.Handle("/api/v1/actors/remove", middleware.MethodMiddleware(
		middleware.AuthorizationMiddleware(
			middleware.PermissionsMiddleware(
				http.HandlerFunc(api.RemoveInfoAboutActor), api.core, variables.AdminRole, api.logger),
			api.core, api.logger),
		http.MethodPost,
		api.logger))

	// Films handlers
	api.mux.Handle("/api/v1/films", middleware.MethodMiddleware(
		middleware.AuthorizationMiddleware(
			http.HandlerFunc(api.GetFilms),
			api.core, api.logger),
		http.MethodGet,
		api.logger))

	api.mux.Handle("/api/v1/films/search", middleware.MethodMiddleware(
		middleware.AuthorizationMiddleware(
			http.HandlerFunc(api.SearchFilms),
			api.core, api.logger),
		http.MethodGet,
		api.logger))

	api.mux.Handle("/api/v1/films/add", middleware.MethodMiddleware(
		middleware.AuthorizationMiddleware(
			middleware.PermissionsMiddleware(
				http.HandlerFunc(api.AddFilm), api.core, variables.AdminRole, api.logger),
			api.core, api.logger),
		http.MethodPost,
		api.logger))

	api.mux.Handle("/api/v1/films/edit", middleware.MethodMiddleware(
		middleware.AuthorizationMiddleware(
			middleware.PermissionsMiddleware(
				http.HandlerFunc(api.EditFilm), api.core, variables.AdminRole, api.logger),
			api.core, api.logger),
		http.MethodPost,
		api.logger))

	api.mux.Handle("/api/v1/films/remove", middleware.MethodMiddleware(
		middleware.AuthorizationMiddleware(
			middleware.PermissionsMiddleware(
				http.HandlerFunc(api.RemoveFilm), api.core, variables.AdminRole, api.logger),
			api.core, api.logger),
		http.MethodPost,
		api.logger))

	return api
}

func (api *API) GetActors(w http.ResponseWriter, r *http.Request) {
	size, page := util.Pagination(r)

	actors, err := api.core.GetActors(uint64((page-1)*size), size)
	if err != nil {
		util.SendResponse(w, r, http.StatusNotFound, nil, variables.ActorsNotFoundError, err, api.logger)
		return
	}
	util.SendResponse(w, r, http.StatusOK, actors, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) GetFilms(w http.ResponseWriter, r *http.Request) {
	sortedBy := r.URL.Query().Get("sort_by")
	pageSize, page := util.Pagination(r)

	films, err := api.core.GetFilms(uint64((page-1)*pageSize), pageSize, sortedBy)
	if err != nil {
		util.SendResponse(w, r, http.StatusNotFound, nil, variables.FilmsNotFoundError, err, api.logger)
		return
	}
	util.SendResponse(w, r, http.StatusOK, films, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) SearchFilms(w http.ResponseWriter, r *http.Request) {
	filmName := r.URL.Query().Get("film_name")
	actorName := r.URL.Query().Get("actor_name")

	film, err := api.core.FindFilm(filmName, actorName)
	if err != nil {
		util.SendResponse(w, r, http.StatusNotFound, nil, variables.FilmNotFoundError, err, api.logger)
		return
	}
	util.SendResponse(w, r, http.StatusOK, film, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) AddInfoAboutActor(w http.ResponseWriter, r *http.Request) {
	var addActorRequest communication.AddActorRequest

	err := util.GetRequestBody(w, r, &addActorRequest, api.logger)
	if err != nil {
		return
	}

	err = api.core.AddActor(addActorRequest.Name, addActorRequest.Gender, addActorRequest.BirthDate)
	if err != nil {
		util.SendResponse(w, r, http.StatusConflict, nil, variables.ActorNotAddedError, err, api.logger)
		return
	}
	util.SendResponse(w, r, http.StatusOK, nil, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) EditInfoAboutActor(w http.ResponseWriter, r *http.Request) {
	var editActorRequest communication.EditActorRequest

	err := util.GetRequestBody(w, r, &editActorRequest, api.logger)
	if err != nil {
		return
	}

	err = api.core.EditActor(editActorRequest.Id, editActorRequest.Name, editActorRequest.Gender, editActorRequest.BirthDate, editActorRequest.Films)
	if err != nil {
		util.SendResponse(w, r, http.StatusConflict, nil, variables.ActorNotEditedError, err, api.logger)
		return
	}
	util.SendResponse(w, r, http.StatusOK, nil, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) RemoveInfoAboutActor(w http.ResponseWriter, r *http.Request) {
	var deleteActorRequest communication.DeleteActorRequest

	err := util.GetRequestBody(w, r, &deleteActorRequest, api.logger)
	if err != nil {
		return
	}

	err = api.core.DeleteActor(deleteActorRequest.Id)
	if err != nil {
		util.SendResponse(w, r, http.StatusConflict, nil, variables.ActorNotDeletedError, err, api.logger)
		return
	}
	util.SendResponse(w, r, http.StatusOK, nil, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) AddFilm(w http.ResponseWriter, r *http.Request) {
	var addFilmRequest communication.AddFilmRequest

	err := util.GetRequestBody(w, r, &addFilmRequest, api.logger)
	if err != nil {
		return
	}

	err = api.core.AddFilm(addFilmRequest.Title, addFilmRequest.Description, addFilmRequest.Rating, addFilmRequest.ReleaseDate, addFilmRequest.Crew)
	if err != nil {
		util.SendResponse(w, r, http.StatusConflict, nil, variables.FilmNotAddedError, err, api.logger)
		return
	}
	util.SendResponse(w, r, http.StatusOK, nil, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) EditFilm(w http.ResponseWriter, r *http.Request) {
	var editFilmRequest communication.EditFilmRequest

	err := util.GetRequestBody(w, r, &editFilmRequest, api.logger)
	if err != nil {
		return
	}

	err = api.core.EditFilm(editFilmRequest.Id, editFilmRequest.Title, editFilmRequest.Description, editFilmRequest.Rating, editFilmRequest.ReleaseDate, editFilmRequest.Crew)
	if err != nil {
		util.SendResponse(w, r, http.StatusConflict, nil, variables.FilmNotEditedError, err, api.logger)
		return
	}
	util.SendResponse(w, r, http.StatusOK, nil, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) RemoveFilm(w http.ResponseWriter, r *http.Request) {
	var deleteFilmRequest communication.DeleteFilmRequest

	err := util.GetRequestBody(w, r, &deleteFilmRequest, api.logger)
	if err != nil {
		return
	}

	err = api.core.DeleteFilm(deleteFilmRequest.Id)
	if err != nil {
		util.SendResponse(w, r, http.StatusConflict, nil, variables.FilmNotDeletedError, err, api.logger)
		return
	}
	util.SendResponse(w, r, http.StatusOK, nil, variables.StatusOkMessage, nil, api.logger)
}
