package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arizanovj/courses/auth"
	"github.com/arizanovj/courses/env"
	"github.com/arizanovj/courses/model"
)

type Auth struct {
	Env *env.Env
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {

	login := model.Login{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&login)
	response := &Response{W: w}
	if err != nil {
		response.Err = err
		response.Code = 300
		response.Json()
		return
	}

	login.DB = a.Env.DB

	err = login.Validate()

	if err != nil {
		response.Err = err
		response.Code = 400
		response.Json()
		return
	}

	if err, _ := login.Login(); err != nil {
		fmt.Printf("%+v\n", err)
		loginError := make(map[string]interface{})
		loginError["email|password"] = "Wrong username or password"
		response.Err = loginError
		response.Code = 400
		response.Json()
		return
	}

	j := &auth.Jwt{
		PublicKeyPath:  "/home/jovan/.config/courses/jwtRS256.key.pub",
		PrivateKeyPath: "/home/jovan/.config/courses/jwtRS256.key",
	}
	token, _ := j.CreateToken(1)

	response.Code = 200
	response.Data = token
	response.Json()

}

func (a *Auth) Signup() {

}

func (a *Auth) Private(w http.ResponseWriter, r *http.Request) {
	response := &Response{W: w}
	response.Code = 200
	response.Data = "Private"
	response.Json()
}
