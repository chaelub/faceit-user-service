package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chaelub/faceit-user-service/internal/api"
	"github.com/chaelub/faceit-user-service/internal/api/types"
	"github.com/chaelub/faceit-user-service/internal/event"
	"github.com/chaelub/faceit-user-service/internal/models"
	"github.com/chaelub/faceit-user-service/internal/repo/user"
	"github.com/chaelub/faceit-user-service/internal/server"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	handler     *mux.Router
	userRequest = types.UserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "jd",
		Password:  "12345",
		Email:     "test@test.com",
		Country:   "DE",
	}
	eventHandler *TestEventHandler
)

const (
	contentType = "application/json"
)

type (
	userResponse struct {
		Success bool        `json:"success"`
		Error   string      `json:"error,omitempty"`
		Payload models.User `json:"payload"`
	}
	TestEventHandler struct {
		lastMsg []byte
	}
)

func (eh *TestEventHandler) handleEvent(msg []byte) error {
	eh.lastMsg = msg
	return nil
}

func init() {
	userRepo := user.NewInMemoryUserRepo()
	eventService := event.NewStubEventService()
	eventHandler = new(TestEventHandler)

	eventService.RegisterHandler(models.UserUpdated, eventHandler.handleEvent)

	router := api.NewApiRouter(userRepo, eventService)
	handler = server.Server(router)
}

func Test_AddUser(t *testing.T) {
	srv := httptest.NewServer(handler)
	defer srv.Close()

	address := fmt.Sprintf("%s/user", srv.URL)

	userResp, err := addUser(address)
	require.Nil(t, err)

	require.True(t, userResp.Success)
	require.Empty(t, userResp.Error)
	require.Equal(t, userRequest.FirstName, userResp.Payload.FirstName)
	// todo: check other fields
	require.Empty(t, userResp.Payload.Password)
}

func Test_Events(t *testing.T) {
	srv := httptest.NewServer(handler)
	defer srv.Close()

	address := fmt.Sprintf("%s/user", srv.URL)

	userResp, err := addUser(address)
	require.Nil(t, err)
	require.True(t, userResp.Success)
	require.NotEmpty(t, userResp.Payload.Id)

	address = fmt.Sprintf("%s/user/%d", srv.URL, userResp.Payload.Id)

	updatedUser := userRequest
	updatedUser.Country = "UK"

	body, err := json.Marshal(updatedUser)
	require.Nil(t, err)
	buf := &bytes.Buffer{}
	_, err = buf.Write(body)
	require.Nil(t, err)

	require.Empty(t, eventHandler.lastMsg)

	req, err := http.NewRequest(http.MethodPut, address, buf)
	require.Nil(t, err)
	req.Header.Set("Content-Type", contentType)
	client := http.Client{}
	resp, err := client.Do(req)

	userResp = userResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	require.Nil(t, err)
	require.True(t, userResp.Success)
	require.Equal(t, userResp.Payload.Country, "UK")

	require.NotEmpty(t, eventHandler.lastMsg)
}

func addUser(address string) (userResponse, error) {
	body, err := json.Marshal(userRequest)
	if err != nil {
		return userResponse{}, err
	}
	buf := &bytes.Buffer{}
	_, err = buf.Write(body)
	if err != nil {
		return userResponse{}, err
	}

	resp, err := http.Post(address, contentType, buf)
	if err != nil {
		return userResponse{}, err
	}

	userResp := userResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	if err != nil {
		return userResponse{}, err
	}

	return userResp, nil
}
