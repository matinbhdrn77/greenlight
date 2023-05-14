package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheack", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v/movies/:id", app.showMovieHandler)

	return router
}
