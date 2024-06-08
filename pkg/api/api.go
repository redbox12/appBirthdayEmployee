package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"taskRumbler/pkg/repository"
)

type api struct {
	r  *mux.Router
	db *repository.PGRepo
}

func New(router *mux.Router, db *repository.PGRepo) *api {
	return &api{r: router, db: db}
}

func (api *api) Handle() {
	api.r.HandleFunc("/api/user/new", api.createUser).Methods(http.MethodPost)
	api.r.HandleFunc("/api/user/login", api.authenticate).Methods(http.MethodPost)
	api.r.HandleFunc("/api/users", api.getUsers).Methods(http.MethodGet)

	api.r.Use(api.jwtAuthentication)
}

func (api *api) ListenAndServe(adrr string) error {
	return http.ListenAndServe(adrr, api.r)
}
