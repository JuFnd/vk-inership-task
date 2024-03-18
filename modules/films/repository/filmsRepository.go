package repository

import (
	"database/sql"
	"filmoteka/pkg/models"
	communication "filmoteka/pkg/requests"
	"filmoteka/pkg/variables"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

type FilmRepository struct {
	db *sql.DB
}

//go:generate mockgen -source=core.go -destination=../mocks/core_mock.go -package=mocks

func GetFilmRepository(configDatabase variables.RelationalDataBaseConfig, logger *slog.Logger) (*FilmRepository, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password= %s host=%s port=%d sslmode=%s",
		configDatabase.User, configDatabase.DbName, configDatabase.Password, configDatabase.Host, configDatabase.Port, configDatabase.Sslmode)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error(variables.SqlOpenError, err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error(variables.SqlPingError, err.Error())
		return nil, err
	}

	db.SetMaxOpenConns(configDatabase.MaxOpenConns)

	filmRepository := &FilmRepository{db: db}

	errs := make(chan error)
	go func() {
		errs <- filmRepository.pingDb(configDatabase.Timer, logger)
	}()

	if err := <-errs; err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return filmRepository, nil
}

func (repository *FilmRepository) pingDb(timer uint32, logger *slog.Logger) error {
	var err error
	var retries int

	for retries < variables.MaxRetries {
		err = repository.db.Ping()
		if err == nil {
			return nil
		}

		retries++
		logger.Error(variables.SqlPingError, err.Error())
		time.Sleep(time.Duration(timer) * time.Second)
	}

	logger.Error(variables.SqlMaxPingRetriesError, err)
	return fmt.Errorf(variables.SqlMaxPingRetriesError, err.Error())
}

func (repository *FilmRepository) GetFilms(begin uint64, end uint64, sortType string) (communication.FilmsListResponse, error) {
	var films []models.FilmItem

	var query string
	switch sortType {
	case "name":
		query = "SELECT f.id, f.name, f.description, f.rating, f.releaseDate, a.id, a.name, a.gender, a.birthdate FROM film f JOIN film_actor fa ON f.id = fa.film_id JOIN actor a ON fa.actor_id = a.id ORDER BY f.name LIMIT $1 OFFSET $2"
	case "rating":
		query = "SELECT f.id, f.name, f.description, f.rating, f.releaseDate, a.id, a.name, a.gender, a.birthdate FROM film f JOIN film_actor fa ON f.id = fa.film_id JOIN actor a ON fa.actor_id = a.id ORDER BY f.rating DESC LIMIT $1 OFFSET $2"
	case "release_date":
		query = "SELECT f.id, f.name, f.description, f.rating, f.releaseDate, a.id, a.name, a.gender, a.birthdate FROM film f JOIN film_actor fa ON f.id = fa.film_id JOIN actor a ON fa.actor_id = a.id ORDER BY f.releaseDate LIMIT $1 OFFSET $2"
	default:
		query = "SELECT f.id, f.name, f.description, f.rating, f.releaseDate, a.id, a.name, a.gender, a.birthdate FROM film f JOIN film_actor fa ON f.id = fa.film_id JOIN actor a ON fa.actor_id = a.id ORDER BY f.rating DESC LIMIT $1 OFFSET $2"
	}

	rows, err := repository.db.Query(query, end-begin, begin)
	if err != nil {
		return communication.FilmsListResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var film models.FilmItem
		var actor models.ActorItem
		err := rows.Scan(&film.Id, &film.Title, &film.Description, &film.Rating, &film.ReleaseDate, &actor.Id, &actor.Name, &actor.Gender, &actor.BirthDate)
		if err != nil {
			return communication.FilmsListResponse{}, err
		}

		var existingFilm *models.FilmItem
		for i := range films {
			if films[i].Id == film.Id {
				existingFilm = &films[i]
				break
			}
		}

		if existingFilm == nil {
			film.Crew = []models.ActorItem{actor}
			films = append(films, film)
		} else {
			existingFilm.Crew = append(existingFilm.Crew, actor)
		}
	}

	response := communication.FilmsListResponse{
		Films: films,
	}

	return response, nil
}

func (repository *FilmRepository) FindFilm(filmName string, actorName string) (communication.FindFilmResponse, error) {
	var response communication.FindFilmResponse

	query := `SELECT f.id, f.name AS title, f.description, f.rating, f.releaseDate, a.id AS actor_id, a.name AS actor_name, a.gender, a.birthdate
              FROM film f
              JOIN film_actor fa ON f.id = fa.film_id
              JOIN actor a ON a.id = fa.actor_id
              WHERE f.name ILIKE $1 || '%'
                  OR f.name ILIKE '%' || $1 || '%'
                  OR f.name ILIKE '%' || $1 || ''
                  OR a.name ILIKE $2 || '%'
                  OR a.name ILIKE '%' || $2 || '%'
                  OR a.name ILIKE '%' || $2 || ''
              ORDER BY 
              (CASE 
                  WHEN f.name ILIKE $1 || '%' THEN 1
                  WHEN f.name ILIKE '%' || $1 || '%' THEN 2
                  WHEN f.name ILIKE '%' || $1 || '' THEN 3
                  WHEN a.name ILIKE $2 || '%' THEN 4
                  WHEN a.name ILIKE '%' || $2 || '%' THEN 5
                  WHEN a.name ILIKE '%' || $2 || '' THEN 6
                  ELSE 7
              END)`

	rows, err := repository.db.Query(query, filmName, actorName)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	films := make(map[int]models.FilmShortItem)
	for rows.Next() {
		var filmID int
		var filmTitle string
		var filmDescription string
		var filmRating float64
		var filmReleaseDate string
		var actorID int
		var actorName string
		var actorGender string
		var actorBirthDate string

		err := rows.Scan(&filmID, &filmTitle, &filmDescription, &filmRating, &filmReleaseDate, &actorID, &actorName, &actorGender, &actorBirthDate)
		if err != nil {
			return response, err
		}

		film, ok := films[filmID]
		if !ok {
			film = models.FilmShortItem{
				Id:          filmID,
				Title:       filmTitle,
				Description: filmDescription,
				Rating:      filmRating,
				ReleaseDate: filmReleaseDate,
			}
		}

		actor := models.ActorShortItem{
			Id:        actorID,
			Name:      actorName,
			Gender:    actorGender,
			BirthDate: actorBirthDate,
		}

		film.Crew = append(film.Crew, actor)
		films[filmID] = film
	}

	for _, film := range films {
		response.Films = append(response.Films, film)
	}

	return response, nil
}

func (repository *FilmRepository) AddFilm(title string, description string, rating float64, releaseDate string, crew []int64) error {
	filmQuery := `INSERT INTO film (name, description, rating, releaseDate) VALUES ($1, $2, $3 ,$4) RETURNING id`
	var filmId int
	err := repository.db.QueryRow(filmQuery, title, description, rating, releaseDate).Scan(&filmId)
	if err != nil {
		return err
	}

	for _, actorId := range crew {
		_, err := repository.db.Exec(`INSERT INTO film_actor (film_id, actor_id) VALUES ($1, $2)`, filmId, actorId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository *FilmRepository) EditFilm(id int64, title string, description string, rating float64, releaseDate string, crew []int64) error {
	_, err := repository.db.Exec(`
    UPDATE film
    SET name = COALESCE($1, name),
        description = COALESCE($2, description),
        releaseDate = COALESCE($3, releaseDate)
    WHERE id = $4`,
		title, description, releaseDate, id)
	if err != nil {
		return err
	}

	_, err = repository.db.Exec(`DELETE FROM film_actor WHERE film_id = $1`, id)
	if err != nil {
		return err
	}

	for _, actorId := range crew {
		_, err := repository.db.Exec(`INSERT INTO film_actor (film_id, actor_id) VALUES ($1, $2)`, id, actorId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository *FilmRepository) GetActors(begin uint64, end uint64) (communication.ActorsListResponse, error) {
	actorsMap := make(map[int]models.ActorItem)
	filmsMap := make(map[int][]models.FilmShortItem)

	rows, err := repository.db.Query(`
        SELECT actor.id, actor.name, actor.gender, actor.birthdate,
               film.id, film.name, film.description, film.rating, film.releaseDate
        FROM actor
        JOIN film_actor ON film_actor.actor_id = actor.id
        JOIN film ON film_actor.film_id = film.id
        ORDER BY actor.id;
    `)
	if err != nil {
		return communication.ActorsListResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var actor models.ActorItem
		var film models.FilmShortItem

		err := rows.Scan(&actor.Id, &actor.Name, &actor.Gender, &actor.BirthDate, &film.Id, &film.Title, &film.Description, &film.Rating, &film.ReleaseDate)
		if err != nil {
			return communication.ActorsListResponse{}, err
		}

		actorsMap[actor.Id] = actor

		filmsMap[actor.Id] = append(filmsMap[actor.Id], film)
	}

	for id, actor := range actorsMap {
		actor.Films = filmsMap[id]
		actorsMap[id] = actor
	}

	var actorsList []models.ActorItem
	for _, actor := range actorsMap {
		actorsList = append(actorsList, actor)
	}

	err = rows.Err()
	if err != nil {
		return communication.ActorsListResponse{}, err
	}

	return communication.ActorsListResponse{Actors: actorsList}, nil
}

func (repository *FilmRepository) AddActor(name string, gender string, birthdate string) error {
	actorQuery := `INSERT INTO actor (name, gender, birthdate) VALUES ($1, $2, $3)`
	_, err := repository.db.Exec(actorQuery, name, gender, birthdate)
	if err != nil {
		return err
	}

	return nil
}

func (repository *FilmRepository) EditActor(id int64, name string, gender string, birthdate string, films []int64) error {
	_, err := repository.db.Exec(`
    UPDATE actor
    SET name = COALESCE($1, name),
        gender = COALESCE($2, gender),
        birthdate = COALESCE($3, birthdate)
    WHERE id = $4`,
		name, gender, birthdate, id)
	if err != nil {
		return err
	}

	_, err = repository.db.Exec(`DELETE FROM film_actor WHERE actor_id = $1`, id)
	if err != nil {
		return err
	}

	for _, filmId := range films {
		_, err := repository.db.Exec(`INSERT INTO film_actor (film_id, actor_id) VALUES ($1, $2)`, filmId, id)
		if err != nil {
			return err
		}
	}

	return nil

}

func (repository *FilmRepository) DeleteActor(id int64) error {
	_, err := repository.db.Exec(`DELETE FROM film_actor WHERE actor_id = $1`, id)
	if err != nil {
		return err
	}

	_, err = repository.db.Exec(`DELETE FROM actor WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (repository *FilmRepository) DeleteFilm(id int64) error {
	_, err := repository.db.Exec(`DELETE FROM film_actor WHERE film_id = $1`, id)
	if err != nil {
		return err
	}

	_, err = repository.db.Exec(`DELETE FROM film WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
