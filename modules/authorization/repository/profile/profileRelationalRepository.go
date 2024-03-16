package profile

import (
	"database/sql"
	"errors"
	"filmoteka/pkg/models"
	"filmoteka/pkg/variables"
	"fmt"
	"log/slog"
	"time"
)

type IProfileRelationalRepository interface {
	CreateUser(login string, password string) error
	FindUser(login string) (bool, error)
	GetUser(login string, password string) (*models.UserItem, bool, error)
	GetUserProfileId(login string) (int64, error)
	GetUserRole(login string) (string, error)
}

type ProfileRelationalRepository struct {
	db *sql.DB
}

func GetProfileRepository(configDatabase *variables.RelationalDataBaseConfig, logger *slog.Logger) (*ProfileRelationalRepository, error) {
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

	profileDb := ProfileRelationalRepository{
		db: db,
	}

	errs := make(chan error)
	go func() {
		errs <- profileDb.pingDb(configDatabase.Timer, logger)
	}()

	if err := <-errs; err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return &profileDb, nil
}

func (profileRelationalRepository *ProfileRelationalRepository) pingDb(timer uint32, logger *slog.Logger) error {
	var err error
	var retries int

	for retries < variables.MaxRetries {
		err = profileRelationalRepository.db.Ping()
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

func (profileRelationalRepository *ProfileRelationalRepository) CreateUser(login string, password string) error {
	_, err := profileRelationalRepository.db.Exec(
		"INSERT INTO profile (login, id_password)"+
			"VALUES ($1, (SELECT id FROM password WHERE value = $2))", login, password)
	if err != nil {
		return fmt.Errorf(variables.SqlProfileCreateError, " %w", err)
	}
	return nil
}

func (profileRelationalRepository *ProfileRelationalRepository) FindUser(login string) (bool, error) {
	userItem := &models.UserItem{}

	err := profileRelationalRepository.db.QueryRow(
		"SELECT login FROM profile"+
			"WHERE login = $1", login).Scan(&userItem.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf(variables.ProfileNotFoundError, ": %w", err)
	}
	return true, nil
}

func (profileRelationalRepository *ProfileRelationalRepository) GetUser(login string, password string) (*models.UserItem, bool, error) {
	userItem := &models.UserItem{}

	err := profileRelationalRepository.db.QueryRow(
		"SELECT login FROM profile"+
			"JOIN password ON profile.id_password = password.id"+
			"WHERE profile.login = $1 AND password.value = $2", login, password).Scan(&userItem.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, fmt.Errorf(variables.InvalidEmailOrPasswordError, ": %w", err)
		}
		return nil, false, fmt.Errorf(variables.ProfileNotFoundError, ": %w", err)
	}

	return userItem, true, nil
}

func (profileRelationalRepository *ProfileRelationalRepository) GetUserProfileId(login string) (int64, error) {
	var userId int64

	err := profileRelationalRepository.db.QueryRow("SELECT id FROM profile WHERE login = $1", login).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf(variables.ProfileIdNotFoundByLoginError, " %s", login)
		}
		return 0, fmt.Errorf(variables.FindProfileIdByLoginError, " %w", err)
	}
	return userId, nil
}

func (profileRelationalRepository *ProfileRelationalRepository) GetUserRole(login string) (string, error) {
	var role string

	err := profileRelationalRepository.db.QueryRow("SELECT role.value FROM profile"+
		"JOIN role ON profile.id_role = role.id WHERE profile.login = $1", login).Scan(&role)
	if err != nil {
		return "", fmt.Errorf(variables.ProfileRoleNotFoundByLoginError, " %w", err)
	}

	return role, nil
}
