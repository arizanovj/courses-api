package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/arizanovj/courses/auth"
	"github.com/arizanovj/courses/env"
	"github.com/gorilla/mux"

	"github.com/arizanovj/courses/handler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
	goqu "gopkg.in/doug-martin/goqu.v4"
	_ "gopkg.in/doug-martin/goqu.v4/adapters/mysql"
)

func main() {

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	configDir := usr.HomeDir + "/.config/courses/"
	configName := "courses"

	viper.SetConfigName(configName) // name of config file (without extension)
	viper.AddConfigPath(configDir)  // path to look for the config file in

	err = viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
	//check congfigs
	if !viper.IsSet("db.user") {
		log.Fatal("missing db user")
	}
	if !viper.IsSet("db.password") {
		log.Fatal("missing db password")
	}
	if !viper.IsSet("db.hostname") {
		log.Fatal("missing db hostname")
	}
	if !viper.IsSet("db.port") {
		log.Fatal("missing port number")
	}
	if !viper.IsSet("db.instance") {
		log.Fatal("missing db instance")
	}
	if !viper.IsSet("db.dialect") {
		log.Fatal("missing db dialect")
	}
	if !viper.IsSet("jwt.publicKey") {
		log.Fatal("missing jwt public key")
	}

	dbUser := viper.GetString("db.user")
	dbPassword := viper.GetString("db.password")
	dbHostname := viper.GetString("db.hostname")
	dbPort := viper.GetString("db.port")
	dbInstance := viper.GetString("db.instance")
	dbDialect := viper.GetString("db.dialect")
	jwtPubkey := viper.GetString("jwt.publicKey")

	db, err := sql.Open(dbDialect, dbUser+":"+dbPassword+"@tcp("+dbHostname+":"+dbPort+")/"+dbInstance+"?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// bytes, err := bcrypt.GenerateFromPassword([]byte("111223344"), 14)
	// fmt.Printf("%+v\n", string(bytes))
	qb := goqu.New(dbDialect, db)
	env := env.Env{
		DB:       db,
		QB:       qb,
		BaseDir:  wd,
		AppURL:   "http://localhost:9001",
		ImageDir: "/static/image/",
		VideoDir: "/static/video/",
	}
	j := &auth.Jwt{
		PublicKeyPath: configDir + jwtPubkey,
	}
	authHandler := &handler.Auth{Env: &env}

	courseHandle := &handler.Course{Env: &env}
	videoHandle := &handler.Video{Env: &env}
	usersHandle := &handler.User{Env: &env}
	resp := &handler.Response{}
	r := mux.NewRouter().PathPrefix("v1").Subrouter()
	r.Handle("/auth/login/", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.Wrap(http.HandlerFunc(authHandler.Login)),
	))
	r.Handle("/auth/private/", negroni.New(
		negroni.HandlerFunc(j.Validate),
		negroni.Wrap(http.HandlerFunc(authHandler.Private)),
	))
	r.Handle("/courses/", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.Wrap(http.HandlerFunc(courseHandle.All)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/courses/{id}", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.Wrap(http.HandlerFunc(courseHandle.Get)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/courses/{id}", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.Wrap(http.HandlerFunc(courseHandle.Update)),
	)).Methods("PUT", "OPTIONS")

	r.Handle("/courses/{id}", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.Wrap(http.HandlerFunc(courseHandle.Delete)),
	)).Methods("DELETE", "OPTIONS")

	r.Handle("/courses/", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.Wrap(http.HandlerFunc(courseHandle.Create)),
	)).Methods("POST", "OPTIONS")

	r.Handle("/courses/{id}/cover", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.Wrap(http.HandlerFunc(courseHandle.CreateCover)),
	)).Methods("POST", "OPTIONS")

	r.Handle("/courses/{id}/cover", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.Wrap(http.HandlerFunc(courseHandle.UpdateCover)),
	)).Methods("PUT", "OPTIONS")

	r.Handle("/videos/", negroni.New(

		negroni.HandlerFunc(resp.CORS),

		negroni.Wrap(http.HandlerFunc(videoHandle.All)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/videos/{id}", negroni.New(

		negroni.HandlerFunc(resp.CORS),

		negroni.Wrap(http.HandlerFunc(videoHandle.Get)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/videos/{id}", negroni.New(

		negroni.HandlerFunc(resp.CORS),

		negroni.Wrap(http.HandlerFunc(videoHandle.Update)),
	)).Methods("PUT", "OPTIONS")

	r.Handle("/videos/{id}", negroni.New(

		negroni.HandlerFunc(resp.CORS),

		negroni.Wrap(http.HandlerFunc(videoHandle.Delete)),
	)).Methods("DELETE", "OPTIONS")

	r.Handle("/videos/", negroni.New(

		negroni.HandlerFunc(resp.CORS),

		negroni.Wrap(http.HandlerFunc(videoHandle.Create)),
	)).Methods("POST", "OPTIONS")

	r.Handle("/videos/{id}/cover", negroni.New(

		negroni.HandlerFunc(resp.CORS),

		negroni.Wrap(http.HandlerFunc(videoHandle.CreateCover)),
	)).Methods("POST", "OPTIONS")

	r.Handle("/videos/{id}/cover", negroni.New(

		negroni.HandlerFunc(resp.CORS),

		negroni.Wrap(http.HandlerFunc(videoHandle.UpdateCover)),
	)).Methods("PUT", "OPTIONS")

	r.Handle("/videos/{id}/src", negroni.New(

		negroni.HandlerFunc(resp.CORS),

		negroni.Wrap(http.HandlerFunc(videoHandle.CreateSrc)),
	)).Methods("POST", "OPTIONS")

	r.Handle("/videos/{id}/src", negroni.New(

		negroni.HandlerFunc(resp.CORS),

		negroni.Wrap(http.HandlerFunc(videoHandle.UpdateSrc)),
	)).Methods("PUT", "OPTIONS")

	r.Handle("/users/", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.Wrap(http.HandlerFunc(usersHandle.All)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/users/", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.Wrap(http.HandlerFunc(usersHandle.Create)),
	)).Methods("POST", "OPTIONS")

	r.Handle("/users/{id}", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.Wrap(http.HandlerFunc(usersHandle.Get)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/users/{id}", negroni.New(
		negroni.HandlerFunc(resp.CORS),
		negroni.Wrap(http.HandlerFunc(usersHandle.Update)),
	)).Methods("PUT", "OPTIONS")

	r.Handle("/users/{id}", negroni.New(

		negroni.HandlerFunc(resp.CORS),

		negroni.Wrap(http.HandlerFunc(usersHandle.Delete)),
	)).Methods("DELETE", "OPTIONS")

	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.Handle("/", r)

	http.ListenAndServe(":9001", nil)
}
