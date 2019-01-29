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

var validImageFileTypes = map[string]string{
	"png":  "image/png",
	"jpeg": "image/jpeg",
	"jpg":  "image/jpeg",
}

var validVideoFileTypes = map[string]string{
	"webm": "video/webm",
	"mp4":  "ivideo/mp4",
}

type Video struct {
	Env *env.Env
}

func (a *Video) All(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}
	video := &model.Video{Env: a.Env}
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
		Model: model.Video{},
	}
	filter.SetFilterParams(r.URL.Query())

	courses, err := video.Get(paginator, filter)

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

func (a *Video) Create(w http.ResponseWriter, r *http.Request) {

	response := &Response{W: w}
	video := &model.Video{Env: a.Env}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&video)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	//err = video.Validate()

	lastID, err := video.Create()

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

func (a *Video) Get(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}
	vars := mux.Vars(r)

	ID, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	video := &model.Video{Env: a.Env}
	videoData, err := video.GetByID(ID)
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	if videoData.Cover != nil {
		coverPath := a.Env.AppURL + a.Env.ImageDir + *(videoData.Cover)
		video.Cover = &coverPath
	}

	if videoData.Src != nil {
		videoPath := a.Env.AppURL + a.Env.VideoDir + *(videoData.Src)
		video.Src = &videoPath
	}

	response.Code = 200
	response.Data = videoData
	response.Json()

}
func (a *Video) Delete(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}
	vars := mux.Vars(r)

	ID, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	video := &model.Video{Env: a.Env, ID: ID}
	err = video.Delete()
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
func (a *Video) Update(w http.ResponseWriter, r *http.Request) {

	response := &Response{W: w}
	video := &model.Video{Env: a.Env}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&video)

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	err = video.Update()
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	response.Code = 200
	response.Data = &video.ID
	response.Json()

}

func (a *Video) CreateCover(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}

	vars := mux.Vars(r)

	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	video := &model.Video{Env: a.Env, ID: ID}
	video, err = video.GetByID(ID)

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

	fileLib := model.File{
		File:       file,
		Header:     header,
		Prefix:     "video_cover_",
		ValidTypes: validImageFileTypes,
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

	video.Cover = &image

	err = video.UpdateCover()

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()

	}

	response.Code = 200
	response.Data = a.Env.AppURL + a.Env.ImageDir + image
	response.Json()

}

func (a *Video) UpdateCover(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}

	vars := mux.Vars(r)

	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	video := &model.Video{Env: a.Env, ID: ID}
	video, err = video.GetByID(ID)

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

	fileLib := model.File{
		File:       file,
		Header:     header,
		Prefix:     "video_cover_",
		ValidTypes: validImageFileTypes,
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
	path := a.Env.BaseDir + a.Env.ImageDir + *(video.Cover)
	if _, err := os.Stat(path); err == nil {
		err = os.Remove(path)
	}

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	video.Cover = &image

	err = video.UpdateCover()

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()

	}

	response.Code = 200
	response.Data = a.Env.AppURL + a.Env.ImageDir + image
	response.Json()

}

func (a *Video) CreateSrc(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}

	vars := mux.Vars(r)

	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	video := &model.Video{Env: a.Env, ID: ID}
	video, err = video.GetByID(ID)

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

	fileLib := model.File{
		File:       file,
		Header:     header,
		Prefix:     "video_src_",
		ValidTypes: validVideoFileTypes,
		Path:       a.Env.BaseDir + a.Env.VideoDir,
	}
	err = fileLib.Validate()
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	videoPath, err := fileLib.SaveFile()

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	video.Src = &videoPath

	err = video.UpdateSrc()

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()

	}

	response.Code = 200
	response.Data = a.Env.AppURL + a.Env.VideoDir + videoPath
	response.Json()

}

func (a *Video) UpdateSrc(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}

	vars := mux.Vars(r)

	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	video := &model.Video{Env: a.Env, ID: ID}
	video, err = video.GetByID(ID)

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

	fileLib := model.File{
		File:       file,
		Header:     header,
		Prefix:     "video_src_",
		ValidTypes: validVideoFileTypes,
		Path:       a.Env.BaseDir + a.Env.VideoDir,
	}

	err = fileLib.Validate()
	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	videoPath, err := fileLib.SaveFile()

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}
	path := a.Env.BaseDir + a.Env.VideoDir + *(video.Src)

	if _, err := os.Stat(path); err == nil {
		err = os.Remove(path)
	}

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()
		return
	}

	video.Src = &videoPath

	err = video.UpdateSrc()

	if err != nil {
		response.Err = err.Error()
		response.Code = 400
		response.Json()

	}

	response.Code = 200
	response.Data = a.Env.AppURL + a.Env.VideoDir + videoPath
	response.Json()

}
