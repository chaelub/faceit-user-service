package api

import (
	"encoding/json"
	"github.com/chaelub/faceit-user-service/internal/api/types"
	"github.com/chaelub/faceit-user-service/internal/models"
	"log"
	"net/http"
	"os"
)

type (
	userProviderI interface {
		New(models.User) (models.User, error)
		Get(int64) (models.User, error)
		Update(models.User) error
		Delete(int64) error
		Find(map[string][]string) ([]models.User, error)
	}

	eventServiceI interface {
		Notify(models.Event) error
	}

	API struct {
		userProvider userProviderI
		eventService eventServiceI
		log          *log.Logger
	}
)

func (a *API) Success(w http.ResponseWriter, payload interface{}) error {
	resp := types.Response{
		Success: true,
		Payload: payload,
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
	return nil
}

func (a *API) Failure(w http.ResponseWriter, status int, err error) error {
	resp := types.Response{
		Success: false,
		Error:   err.Error(),
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	w.Write(respBytes)
	return nil
}

// @title FACEIT User service API
// @version 0.1

// @host localhost:3000
// @BasePath /
func NewApiRouter(userProvider userProviderI, events eventServiceI) *API {
	return &API{
		userProvider: userProvider,
		eventService: events,
		log:          log.New(os.Stderr, "[Router] ", log.LstdFlags),
	}
}
