package middleware

import (
	"context"
	"filmoteka/pkg/util"
	"filmoteka/pkg/variables"
	"log/slog"
	"net/http"
	"strconv"
)

type ICore interface {
	GetUserId(ctx context.Context, sid string) (int, error)
}

func MethodMiddleware(next http.Handler, method string, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			util.SendResponse(w, r, http.StatusMethodNotAllowed, nil, variables.StatusMethodNotAllowedError, nil, logger)
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthorizationMiddleware(next http.Handler, core ICore, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if err != nil {
			logger.Error(variables.SessionNotFoundError, r.Method, strconv.Itoa(http.StatusUnauthorized), r.URL.Path, err.Error())
			next.ServeHTTP(w, r)
			return
		}

		userId, err := core.GetUserId(r.Context(), session.Value)
		if err != nil || userId == 0 {
			logger.Error(variables.UserNotAuthorized, r.Method, strconv.Itoa(http.StatusUnauthorized), r.URL.Path, err.Error())
			next.ServeHTTP(w, r)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), variables.UserIDKey, userId))
		r = r.WithContext(context.WithValue(r.Context(), variables.SessionIDKey, session))
		next.ServeHTTP(w, r)
	})
}
