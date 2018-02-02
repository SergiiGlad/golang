package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-team-room/conf"
	"go-team-room/controllers"
	"go-team-room/models/dto"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Get instance of logger (Formatter, Hook，Level, Output ).
// If you want to use only your log message  It will need use own call logs example
var log = conf.GetLog()

var i = 0

func createProfile(service controllers.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqUserDto, err := userDtoFromReq(r)

		if err != nil {
			responseError(w, err)
			return
		}

		respUserDto, err := service.CreateUser(&reqUserDto)

		if err != nil {
			responseError(w, err)
			return
		}

		respBody, err := json.Marshal(respUserDto)
		_, err = w.Write(respBody)

		if err != nil {
			responseError(w, err)
			return
		}
	}
}

func updateProfile(service controllers.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userDto, err := userDtoFromReq(r)

		if err != nil {
			responseError(w, err)
			return
		}

		idStr := mux.Vars(r)["user_id"]
		id, err := strconv.Atoi(idStr)

		if err != nil {
			responseError(w, err)
			return
		}

		respUserDto, err := service.UpdateUser(int64(id), &userDto)

		if err != nil {
			responseError(w, err)
			return
		}

		respBody, err := json.Marshal(respUserDto)
		_, err = w.Write(respBody)

		if err != nil {
			responseError(w, err)
			return
		}
	}
}

func deleteProfile(service controllers.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["user_id"]
		id, err := strconv.Atoi(idStr)

		if err != nil {
			responseError(w, err)
			return
		}

		respUserDto, err := service.DeleteUser(int64(id))

		if err != nil {
			responseError(w, err)
			return
		}

		respBody, err := json.Marshal(respUserDto)
		_, err = w.Write(respBody)

		if err != nil {
			responseError(w, err)
			return
		}
	}
}

func responseError(w http.ResponseWriter, err error) {
	rerror := dto.ResponseError{err.Error()}
	log.Println(err)

	body, err := json.Marshal(rerror)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err)
	}

	http.Error(w, string(body), http.StatusBadRequest)
}

func userDtoFromReq(request *http.Request) (dto.RequestUserDto, error) {
	body, err := ioutil.ReadAll(request.Body)
	userDto := dto.RequestUserDto{}

	if err != nil {
		return userDto, err
	}

	err = json.Unmarshal(body, &userDto)

	if err != nil {
		return userDto, err
	}

	return userDto, err
}
