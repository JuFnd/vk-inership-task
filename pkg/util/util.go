package util

import (
	"encoding/json"
	"filmoteka/pkg/variables"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func SendResponse(w http.ResponseWriter, r *http.Request, status int, body any, errorMessage string, handlerError error, logger *slog.Logger) {
	jsonResponse, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(variables.JsonPackFailedError, r.Method, strconv.Itoa(http.StatusInternalServerError), r.URL.Path, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(jsonResponse)
	if err != nil {
		logger.Error(variables.ResponseSendFailedError, r.Method, strconv.Itoa(http.StatusInternalServerError), r.URL.Path, err.Error())
		return
	}

	if handlerError != nil {
		logger.Error(errorMessage, r.Method, strconv.Itoa(status), r.URL.Path, handlerError.Error())
		return
	}

	logger.Error(errorMessage, r.Method, strconv.Itoa(status), r.URL.Path, nil)
}

func GetRequestBody(w http.ResponseWriter, r *http.Request, requestObject any, logger *slog.Logger) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendResponse(w, r, http.StatusBadRequest, nil, variables.StatusBadRequestError, err, logger)
		return err
	}

	err = json.Unmarshal(body, requestObject)
	if err != nil {
		SendResponse(w, r, http.StatusBadRequest, nil, variables.StatusBadRequestError, err, logger)
		return err
	}
	return nil
}

func GetCookie(name string, value string, path string, httpOnly bool, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Expires:  expires,
		HttpOnly: httpOnly,
	}
}

func RandStringRunes(seed int) string {
	symbols := make([]rune, seed)
	for i := range symbols {
		symbols[i] = variables.LetterRunes[rand.Intn(len(variables.LetterRunes))]
	}
	return string(symbols)
}
