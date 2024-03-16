package delivery

import (
	"filmoteka/modules/authorization/usecase"
	"filmoteka/pkg/middleware"
	communication "filmoteka/pkg/requests"
	"filmoteka/pkg/util"
	"filmoteka/pkg/variables"
	"log/slog"
	"net/http"
	"time"
)

type API struct {
	core   usecase.ICore
	logger *slog.Logger
	mux    *http.ServeMux
}

func (api *API) ListenAndServe(appConfig variables.AuthorizationAppConfig) error {
	err := http.ListenAndServe(":"+appConfig.Port, api.mux)
	if err != nil {
		api.logger.Error(variables.ListenAndServeError, err.Error())
		return err
	}
	return nil
}

func GetAuthorizationApi(authCore *usecase.Core, authLogger *slog.Logger) *API {
	api := &API{
		core:   authCore,
		logger: authLogger,
		mux:    http.NewServeMux(),
	}

	// Signin handler
	api.mux.Handle("/signin", middleware.MethodMiddleware(
		http.HandlerFunc(api.Signin),
		http.MethodPost,
		api.logger))

	// Signup handler
	api.mux.Handle("/signup", middleware.MethodMiddleware(
		http.HandlerFunc(api.Signup),
		http.MethodPost,
		api.logger))

	// Logout handler
	api.mux.Handle("/logout", middleware.MethodMiddleware(
		middleware.AuthorizationMiddleware(
			http.HandlerFunc(api.LogoutSession),
			api.core, api.logger),
		http.MethodPost,
		api.logger))

	return api
}

func (api *API) Signin(w http.ResponseWriter, r *http.Request) {
	var signinRequest communication.SigninRequest

	err := util.GetRequestBody(w, r, &signinRequest, api.logger)
	if err != nil {
		return
	}

	user, found, err := api.core.FindUserAccount(signinRequest.Login, signinRequest.Password)
	if err != nil {
		util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.StatusInternalServerError, err, api.logger)
		return
	}

	if !found {
		util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.StatusUnauthorizedError, nil, api.logger)
		return
	}

	session, err := api.core.CreateSession(r.Context(), user.Login)
	if err != nil {
		util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.SessionCreateError, err, api.logger)
		return
	}

	authorizationCookie := util.GetCookie(variables.SessionCookieName, session.SID, "/", variables.HttpOnly, session.ExpiresAt)
	http.SetCookie(w, authorizationCookie)
	util.SendResponse(w, r, http.StatusOK, nil, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) Signup(w http.ResponseWriter, r *http.Request) {
	var signupRequest communication.SignupRequest

	err := util.GetRequestBody(w, r, &signupRequest, api.logger)
	if err != nil {
		return
	}

	found, err := api.core.FindUserByLogin(signupRequest.Login)
	if err != nil {
		util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.StatusInternalServerError, err, api.logger)
		return
	}

	if found {
		util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.UserAlreadyExistsError, nil, api.logger)
		return
	}

	err = api.core.CreateUserAccount(signupRequest.Login, signupRequest.Password)
	if err != nil && err.Error() == variables.InvalidEmailOrPasswordError {
		util.SendResponse(w, r, http.StatusBadRequest, nil, variables.InvalidEmailOrPasswordError, err, api.logger)
		return
	}
	if err != nil {
		util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.UserAlreadyExistsError, err, api.logger)
		return
	}
	util.SendResponse(w, r, http.StatusOK, nil, variables.StatusOkMessage, nil, api.logger)
}

func (api *API) LogoutSession(w http.ResponseWriter, r *http.Request) {
	session, isAuth := r.Context().Value(variables.SessionIDKey).(*http.Cookie)
	if !isAuth {
		util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.SessionNotFoundError, nil, api.logger)
		return
	}

	found, err := api.core.FindActiveSession(r.Context(), session.Value)
	if err != nil {
		util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.StatusUnauthorizedError, err, api.logger)
		return
	}

	if !found {
		util.SendResponse(w, r, http.StatusUnauthorized, nil, variables.SessionNotFoundError, nil, api.logger)
		return
	}

	err = api.core.KillSession(r.Context(), session.Value)
	if err != nil {
		util.SendResponse(w, r, http.StatusInternalServerError, nil, variables.SessionKilledError, err, api.logger)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	util.SendResponse(w, r, http.StatusOK, nil, variables.StatusOkMessage, nil, api.logger)
}
