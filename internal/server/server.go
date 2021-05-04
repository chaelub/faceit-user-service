package server

import (
	_ "github.com/chaelub/faceit-user-service/docs"
	"github.com/chaelub/faceit-user-service/internal/api"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func Server(router *api.API) *mux.Router {
	s := mux.NewRouter()
	s.HandleFunc("/user", router.AddUser).Methods(http.MethodPost)
	s.HandleFunc("/user/{id:[0-9]+}", router.User).Methods(http.MethodGet)
	s.HandleFunc("/user/{id:[0-9]+}", router.DeleteUser).Methods(http.MethodDelete)
	s.HandleFunc("/user/{id:[0-9]+}", router.UpdateUser).Methods(http.MethodPut)
	s.HandleFunc("/user/find", router.FindUsers).Methods(http.MethodGet)

	s.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	return s
}
