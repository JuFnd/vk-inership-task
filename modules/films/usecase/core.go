package usecase

import (
	"context"
	"filmoteka/modules/authorization/proto/authorization"
	communication "filmoteka/pkg/requests"
	"filmoteka/pkg/util"
	"filmoteka/pkg/variables"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type IFilmRepository interface {
	GetFilms(begin uint64, end uint64, sortType string) (communication.FilmsListResponse, error)
	FindFilm(filmName string, actorName string) (communication.FindFilmResponse, error)
	AddFilm(title string, description string, rating float64, releaseDate string, crew []int64) error
	EditFilm(id int64, title string, description string, rating float64, releaseDate string, crew []int64) error
	GetActors(begin uint64, end uint64) (communication.ActorsListResponse, error)
	AddActor(name string, gender string, birthdate string) error
	EditActor(id int64, name string, gender string, birthdate string, films []int64) error
	DeleteActor(id int64) error
	DeleteFilm(id int64) error
}

type Core struct {
	filmRepository IFilmRepository
	client         authorization.AuthorizationClient
	logger         *slog.Logger
}

func GetGrpcClient(port string) (authorization.AuthorizationClient, error) {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf(variables.GrpcConnectError, ": %w", err)
	}
	client := authorization.NewAuthorizationClient(conn)

	return client, nil
}

func GetCore(configGrpc variables.GrpcConfig, films IFilmRepository, logger *slog.Logger) *Core {
	client, err := GetGrpcClient(configGrpc.Port)
	if err != nil {
		logger.Error(variables.GrpcConnectError, ": %w", err)
		return nil
	}
	return &Core{
		filmRepository: films,
		client:         client,
		logger:         logger,
	}
}

func (core *Core) GetFilms(begin uint64, end uint64, sortType string) (communication.FilmsListResponse, error) {
	filmsList, err := core.filmRepository.GetFilms(begin, end, sortType)
	if err != nil {
		core.logger.Error(variables.FilmsListNotFoundError, err)
		return communication.FilmsListResponse{}, err
	}
	return filmsList, nil
}

func (core *Core) FindFilm(filmName string, actorName string) (communication.FindFilmResponse, error) {
	film, err := core.filmRepository.FindFilm(filmName, actorName)
	if err != nil {
		core.logger.Error(variables.FilmNotFoundError, err)
		return communication.FindFilmResponse{}, err
	}
	return film, nil
}

func (core *Core) AddFilm(title string, description string, rating float64, releaseDate string, crew []int64) error {
	if rating < variables.FilmRatingBegin || rating > variables.FilmRatingEnd {
		core.logger.Error(variables.RatingSizeError)
		return fmt.Errorf(variables.RatingSizeError)
	}

	err := util.ValidateStringSize(title, variables.FilmTitleBegin, variables.FilmTitleEnd, variables.TitleSizeError, core.logger)
	if err != nil {
		return err
	}

	err = util.ValidateStringSize(description, variables.FilmDescriptionBegin, variables.FilmDescriptionEnd, variables.DescriptionSizeError, core.logger)
	if err != nil {
		return err
	}

	err = core.filmRepository.AddFilm(title, description, rating, releaseDate, crew)
	if err != nil {
		core.logger.Error(variables.FilmNotAddedError, err)
		return err
	}
	return nil
}

func (core *Core) EditFilm(id int64, title string, description string, rating float64, releaseDate string, crew []int64) error {
	if rating < variables.FilmRatingBegin || rating > variables.FilmRatingEnd {
		core.logger.Error(variables.RatingSizeError)
		return fmt.Errorf(variables.RatingSizeError)
	}

	err := util.ValidateStringSize(title, variables.FilmTitleBegin, variables.FilmTitleEnd, variables.TitleSizeError, core.logger)
	if err != nil {
		return err
	}

	err = util.ValidateStringSize(description, variables.FilmDescriptionBegin, variables.FilmDescriptionEnd, variables.DescriptionSizeError, core.logger)
	if err != nil {
		return err
	}

	err = core.filmRepository.EditFilm(id, title, description, rating, releaseDate, crew)
	if err != nil {
		core.logger.Error(variables.FilmNotEditedError, err)
		return err
	}
	return nil
}

func (core *Core) GetActors(begin uint64, end uint64) (communication.ActorsListResponse, error) {
	actorsList, err := core.filmRepository.GetActors(begin, end)
	if err != nil {
		core.logger.Error(variables.ActorsNotFoundError, err)
		return communication.ActorsListResponse{}, err
	}
	return actorsList, nil
}

func (core *Core) AddActor(name string, gender string, birthdate string) error {
	err := util.ValidateStringSize(name, variables.ActorNameBegin, variables.ActorNameEnd, variables.ActorNameSizeError, core.logger)
	if err != nil {
		return err
	}

	err = core.filmRepository.AddActor(name, gender, birthdate)
	if err != nil {
		core.logger.Error(variables.ActorNotAddedError, err)
		return err
	}
	return nil
}

func (core *Core) EditActor(id int64, name string, gender string, birthdate string, films []int64) error {
	err := util.ValidateStringSize(name, variables.ActorNameBegin, variables.ActorNameEnd, variables.ActorNameSizeError, core.logger)
	if err != nil {
		return err
	}

	err = core.filmRepository.EditActor(id, name, gender, birthdate, films)
	if err != nil {
		core.logger.Error(variables.ActorNotEditedError, err)
		return err
	}
	return nil
}

func (core *Core) DeleteActor(id int64) error {
	err := core.filmRepository.DeleteActor(id)
	if err != nil {
		core.logger.Error(variables.ActorNotDeletedError, err)
		return err
	}
	return nil
}

func (core *Core) DeleteFilm(id int64) error {
	err := core.filmRepository.DeleteFilm(id)
	if err != nil {
		core.logger.Error(variables.FilmNotDeletedError, err)
		return err
	}
	return nil
}

func (core *Core) GetUserRole(ctx context.Context, id int64) (string, error) {
	grpcRequest := authorization.RoleRequest{Id: id}

	grpcResponse, err := core.client.GetRole(ctx, &grpcRequest)
	if err != nil {
		core.logger.Error(variables.GrpcRecievError, err)
		return "", fmt.Errorf(variables.GrpcRecievError, err)
	}
	return grpcResponse.GetRole(), nil
}

func (core *Core) GetUserId(ctx context.Context, sid string) (int64, error) {
	grpcRequest := authorization.FindIdRequest{Sid: sid}

	grpcResponse, err := core.client.GetId(ctx, &grpcRequest)
	if err != nil {
		core.logger.Error(variables.GrpcRecievError, err)
		return 0, fmt.Errorf(variables.GrpcRecievError, err)
	}
	return grpcResponse.Value, nil
}
