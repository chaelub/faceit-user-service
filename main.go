package main

import (
	"github.com/chaelub/faceit-user-service/internal/api"
	"github.com/chaelub/faceit-user-service/internal/event"
	"github.com/chaelub/faceit-user-service/internal/models"
	"github.com/chaelub/faceit-user-service/internal/repo/user"
	"github.com/chaelub/faceit-user-service/internal/server"
	"log"
	"net/http"
	"os"
)

var (
	logger = log.New(os.Stderr, "[Main] ", log.LstdFlags)
)

func main() {

	userRepo := user.NewInMemoryUserRepo()
	eventService := event.NewStubEventService()

	eventService.RegisterHandler(models.UserUpdated, event.NewStubEventHandler(models.UserUpdated))

	router := api.NewApiRouter(userRepo, eventService)
	handler := server.Server(router)

	err := http.ListenAndServe(":3000", handler)
	if err != nil {
		logger.Fatal(err)
	}
}
