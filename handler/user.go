package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/arizanovj/courses/env"
	"github.com/arizanovj/courses/libs"
	"github.com/arizanovj/courses/libs/filter"
	"github.com/arizanovj/courses/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type User struct {
	Env *env.Env
}

func (a *User) All(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}
	user := &model.User{Env: a.Env}
	paginator := &pagination.Paginator{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(paginator, r.URL.Query()); err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	if err := paginator.Validate(); err != nil {
		response.Err = err
		response.Code = 400
		response.Json()
		return
	}
	paginator.Env = a.Env

	filter := &filter.Filter{
		Env:   a.Env,
		Model: model.User{},
	}
	filter.SetFilterParams(r.URL.Query())

	users, err := user.Get(paginator, filter)

	if err != nil {
		response.Err = err
		response.Code = 400
		response.Json()
		return
	}

	response.Code = 200
	response.Data = users
	response.Json()
}

func (a *User) Create(w http.ResponseWriter, r *http.Request) {

	response := &Response{W: w}
	user := &model.User{Env: a.Env}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	lastID, err := user.Create()

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	response.Code = 200
	response.Data = lastID
	response.Json()
}

func (a *User) Get(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}
	vars := mux.Vars(r)

	ID, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	user := &model.User{Env: a.Env}
	userData, err := user.GetByID(ID)
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	response.Code = 200
	response.Data = userData
	response.Json()
}

func (a *User) Update(w http.ResponseWriter, r *http.Request) {

	response := &Response{W: w}
	user := &model.User{Env: a.Env}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	err = user.Update()
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	response.Code = 200
	response.Data = &user.ID
	response.Json()

}

func (a *User) Delete(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}
	vars := mux.Vars(r)

	ID, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	user := &model.User{Env: a.Env, ID: ID}
	err = user.Delete()
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	response.Code = 200
	response.Data = ID
	response.Json()

}
