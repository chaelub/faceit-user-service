package api

import (
	"encoding/json"
	"errors"
	"github.com/chaelub/faceit-user-service/internal/api/types"
	"github.com/chaelub/faceit-user-service/internal/models"
	user2 "github.com/chaelub/faceit-user-service/internal/repo/user"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// AddUser godoc
// @Summary Register new user
// @Produce json
// @Accept  json
// @Param user body types.UserRequest true "New user's data"
// @Success 200 {object} models.publicUser
// @Failure 400,404,500 {object} types.Response
// @Router /user [post]
func (a *API) AddUser(w http.ResponseWriter, req *http.Request) {
	user, err := receiveUser(req)
	if err != nil {
		a.Failure(w, http.StatusBadRequest, err)
		return
	}
	user, err = a.userProvider.New(user)
	if err != nil {
		a.Failure(w, http.StatusInternalServerError, err)
		return
	}
	a.Success(w, user)
}

// User godoc
// @Summary Returns User model by given ID
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} models.publicUser
// @Router /user/{id} [get]
func (a *API) User(w http.ResponseWriter, req *http.Request) {
	userId, err := userId(req)
	if err != nil {
		a.Failure(w, http.StatusBadRequest, err)
		return
	}

	user, err := a.userProvider.Get(userId)
	if err != nil {
		if err == user2.ErrUserNotExists {
			a.Failure(w, http.StatusNotFound, err)
			return
		}
		a.Failure(w, http.StatusInternalServerError, err)
		return
	}
	a.Success(w, user)
}

// UpdateUser godoc
// @Summary Updates User model
// @Produce json
// @Accept  json
// @Param user body types.UserRequest true "New user's data"
// @Param id path integer true "User ID"
// @Success 200 {object} models.publicUser
// @Router /user/{id} [put]
func (a *API) UpdateUser(w http.ResponseWriter, req *http.Request) {
	userId, err := userId(req)
	if err != nil {
		a.Failure(w, http.StatusBadRequest, err)
		return
	}

	user, err := receiveUser(req)
	if err != nil {
		a.Failure(w, http.StatusBadRequest, err)
		return
	}

	user.Id = userId
	err = a.userProvider.Update(user)
	if err != nil {
		a.Failure(w, http.StatusInternalServerError, err)
		return
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		a.Failure(w, http.StatusInternalServerError, err)
		return
	}

	evt := models.Event{
		Type:    models.UserUpdated,
		Message: userBytes,
	}
	if err = a.eventService.Notify(evt); err != nil {
		a.log.Printf("can't send notification about changed user data: %v", err)
	}

	a.Success(w, user)
}

// DeleteUser godoc
// @Summary Deletes whole user data
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} models.publicUser
// @Router /user/{id} [delete]
func (a *API) DeleteUser(w http.ResponseWriter, req *http.Request) {
	userId, err := userId(req)
	if err != nil {
		a.Failure(w, http.StatusBadRequest, err)
		return
	}

	err = a.userProvider.Delete(userId)
	if err != nil {
		a.Failure(w, http.StatusInternalServerError, err)
		return
	}
	a.Success(w, "success")
}

// FindUsers godoc
// @Summary Returns a list of User models founded by given criteria
// @Produce json
// @Param email query string false "Email regexp template"
// @Param country query string false "Country regexp template"
// @Success 200 {array} models.publicUser
// @Router /user/find [get]
func (a *API) FindUsers(w http.ResponseWriter, req *http.Request) {
	queryVals := req.URL.Query()
	users, err := a.userProvider.Find(queryVals)
	if err != nil {
		a.Failure(w, http.StatusInternalServerError, err)
		return
	}
	if len(users) == 0 {
		a.Failure(w, http.StatusNotFound, errors.New("not found"))
		return
	}
	a.Success(w, users)
}

func userId(req *http.Request) (int64, error) {
	vars := mux.Vars(req)
	userIdS, got := vars["id"]
	if !got {
		return 0, errors.New("empty user id")
	}
	userId, err := strconv.ParseInt(userIdS, 10, 64)
	if userId <= 0 {
		return 0, errors.New("bad user id")
	}
	return userId, err
}

func receiveUser(req *http.Request) (models.User, error) {
	userData := types.UserRequest{}
	err := json.NewDecoder(req.Body).Decode(&userData)
	user := models.User{
		Id:        0,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Nickname:  userData.Nickname,
		Password:  userData.Password,
		Email:     userData.Email,
		Country:   userData.Country,
	}
	return user, err
}
