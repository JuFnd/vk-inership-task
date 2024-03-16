package usecase

import (
	"context"
	"filmoteka/modules/authorization/repository/profile"
	"filmoteka/modules/authorization/repository/session"
	"filmoteka/pkg/models"
	"filmoteka/pkg/util"
	"filmoteka/pkg/variables"
	"fmt"
	"log/slog"
	"regexp"
	"sync"
	"time"
)

type ICore interface {
	KillSession(ctx context.Context, sid string) error
	FindActiveSession(ctx context.Context, sid string) (bool, error)
	CreateSession(ctx context.Context, login string) (models.Session, error)
	CreateUserAccount(login string, password string) error
	FindUserByLogin(login string) (bool, error)
	FindUserAccount(login string, password string) (*models.UserItem, bool, error)
	GetUserId(ctx context.Context, sid string) (int64, error)
	GetUserRole(login string) (string, error)
}

type Core struct {
	sessions session.SessionCacheRepository
	logger   *slog.Logger
	mutex    sync.RWMutex
	profiles profile.IProfileRelationalRepository
}

func GetCore(profileConfig *variables.RelationalDataBaseConfig, sessionConfig variables.CacheDataBaseConfig, logger *slog.Logger) (*Core, error) {
	sessionRepository, err := session.GetSessionRepository(sessionConfig, logger)
	if err != nil {
		logger.Error(variables.SessionRepositoryNotActiveError)
		return nil, err
	}

	profileRepository, err := profile.GetProfileRepository(profileConfig, logger)
	if err != nil {
		logger.Error(variables.ProfileRepositoryNotActiveError)
		return nil, err
	}

	core := Core{
		sessions: *sessionRepository,
		logger:   logger.With(variables.ModuleLogger, variables.CoreModuleLogger),
		profiles: profileRepository,
	}

	return &core, nil
}

func (core *Core) CreateSession(ctx context.Context, login string) (models.Session, error) {
	sid := util.RandStringRunes(32)

	newSession := models.Session{
		Login:     login,
		SID:       sid,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	core.mutex.Lock()
	sessionAdded, err := core.sessions.SaveSessionCache(ctx, newSession, core.logger)
	defer core.mutex.Unlock()

	if !sessionAdded && err != nil {
		return models.Session{}, err
	}

	if !sessionAdded {
		return models.Session{}, nil
	}

	return newSession, nil
}

func (core *Core) KillSession(ctx context.Context, sid string) error {
	core.mutex.Lock()
	_, err := core.sessions.DeleteSessionCache(ctx, sid, core.logger)
	defer core.mutex.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (core *Core) FindActiveSession(ctx context.Context, sid string) (bool, error) {
	core.mutex.RLock()
	found, err := core.sessions.GetSessionCache(ctx, sid, core.logger)
	defer core.mutex.RUnlock()

	if err != nil {
		return false, err
	}

	return found, nil
}

func (core *Core) CreateUserAccount(login string, password string) error {
	matched, err := regexp.MatchString(variables.LoginRegexp, login)
	if err != nil {
		core.logger.Error(variables.StatusInternalServerError, err.Error())
		return fmt.Errorf(variables.StatusInternalServerError, " %w", err)
	}
	if !matched {
		core.logger.Error(variables.InvalidEmailOrPasswordError)
		return fmt.Errorf(variables.InvalidEmailOrPasswordError)
	}

	hashPassword := util.HashPassword(password)
	err = core.profiles.CreateUser(login, hashPassword)
	if err != nil {
		core.logger.Error(variables.CreateProfileError, err.Error())
		return err
	}

	return nil
}

func (core *Core) FindUserByLogin(login string) (bool, error) {
	found, err := core.profiles.FindUser(login)
	if err != nil {
		core.logger.Error(variables.ProfileNotFoundError, err.Error())
		return false, err
	}

	return found, nil
}

func (core *Core) FindUserAccount(login string, password string) (*models.UserItem, bool, error) {
	hashPassword := util.HashPassword(password)
	user, found, err := core.profiles.GetUser(login, hashPassword)
	if err != nil {
		core.logger.Error(variables.ProfileNotFoundError, err.Error())
		return nil, false, err
	}
	return user, found, nil
}

func (core *Core) GetUserId(ctx context.Context, sid string) (int64, error) {
	login, err := core.sessions.GetUserLogin(ctx, sid, core.logger)
	if err != nil {
		return 0, err
	}

	id, err := core.profiles.GetUserProfileId(login)
	if err != nil {
		core.logger.Error(variables.GetProfileError, " id: %v", err)
		return 0, err
	}
	return id, nil
}

func (core *Core) GetUserRole(login string) (string, error) {
	role, err := core.profiles.GetUserRole(login)
	if err != nil {
		core.logger.Error(variables.GetProfileRoleError, err.Error())
		return "", fmt.Errorf(variables.GetProfileRoleError, " %w", err)
	}

	return role, nil
}
