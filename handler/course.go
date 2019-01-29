package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/arizanovj/courses/env"
	"github.com/arizanovj/courses/libs"
	"github.com/arizanovj/courses/libs/filter"
	"github.com/arizanovj/courses/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type Course struct {
	Env *env.Env
}

func (a *Course) All(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}
	course := &model.Course{Env: a.Env}
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
		Model: model.Course{},
	}
	filter.SetFilterParams(r.URL.Query())

	courses, err := course.Get(paginator, filter)

	if err != nil {
		response.Err = err
		response.Code = 400
		response.Json()
		return
	}

	response.Code = 200
	response.Data = courses
	response.Json()
}

func (a *Course) Create(w http.ResponseWriter, r *http.Request) {

	response := &Response{W: w}
	course := &model.Course{Env: a.Env}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&course)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	//err = course.Validate()

	lastID, err := course.Create()

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

func (a *Course) Get(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}
	vars := mux.Vars(r)

	ID, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	course := &model.Course{Env: a.Env}
	courseData, err := course.GetByID(ID)
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	path := a.Env.AppURL + a.Env.ImageDir + *(course.Cover)
	course.Cover = &path
	response.Code = 200
	response.Data = courseData
	response.Json()

}
func (a *Course) Delete(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}
	vars := mux.Vars(r)

	ID, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	course := &model.Course{Env: a.Env, ID: ID}
	err = course.Delete()
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
func (a *Course) Update(w http.ResponseWriter, r *http.Request) {

	response := &Response{W: w}
	course := &model.Course{Env: a.Env}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&course)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	err = course.Update()
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	response.Code = 200
	response.Data = &course.ID
	response.Json()

}

func (a *Course) CreateCover(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}

	vars := mux.Vars(r)

	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	course := &model.Course{Env: a.Env, ID: ID}
	course, err = course.GetByID(ID)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	file, header, err := r.FormFile("file")

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	validFileTypes := map[string]string{
		"png":  "image/png",
		"jpeg": "image/jpeg",
		"jpg":  "image/jpeg",
	}
	fileLib := model.File{
		File:       file,
		Header:     header,
		Prefix:     "course_cover_",
		ValidTypes: validFileTypes,
		Path:       a.Env.BaseDir + a.Env.ImageDir,
	}
	err = fileLib.Validate()
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	image, err := fileLib.SaveFile()

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	course.Cover = &image

	err = course.UpdateCover()

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()

	}

	response.Code = 200
	response.Data = a.Env.AppURL + a.Env.ImageDir + image
	response.Json()

}

func (a *Course) UpdateCover(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}

	vars := mux.Vars(r)

	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	course := &model.Course{Env: a.Env, ID: ID}
	course, err = course.GetByID(ID)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	file, header, err := r.FormFile("file")

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	validFileTypes := map[string]string{
		"png":  "image/png",
		"jpeg": "image/jpeg",
		"jpg":  "image/jpeg",
	}
	fileLib := model.File{
		File:       file,
		Header:     header,
		Prefix:     "course_cover_",
		ValidTypes: validFileTypes,
		Path:       a.Env.BaseDir + a.Env.ImageDir,
	}
	err = fileLib.Validate()
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	image, err := fileLib.SaveFile()

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	path := (a.Env.BaseDir + a.Env.ImageDir + *(course.Cover))
	if _, err := os.Stat(path); err == nil {
		err = os.Remove(path)
	}
	err = os.Remove(path)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	course.Cover = &image

	err = course.UpdateCover()

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()

	}

	response.Code = 200
	response.Data = a.Env.AppURL + a.Env.ImageDir + image
	response.Json()

}
