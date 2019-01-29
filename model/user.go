package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/arizanovj/courses/env"

	"github.com/arizanovj/courses/libs"
	"github.com/arizanovj/courses/libs/filter"
	"golang.org/x/crypto/bcrypt"
	goqu "gopkg.in/doug-martin/goqu.v4"
)

type User struct {
	ID           int64    `json:"id" filter:"id,number"`
	Email        string   `json:"email" filter:"email,string"`
	FirstName    string   `json:"first_name" filter:"first_name,string"`
	LastName     string   `json:"last_name" filter:"last_name,string"`
	PasswordHash string   `json:"-" filter:"-"`
	Password     string   `json:"password" filter:"-"`
	IsAdmin      *bool    `json:"is_admin" filter:"is_admin,string"`
	CreatedAt    string   `json:"created_at" filter:"created_at,string"`
	UpdatedAt    string   `json:"updated_at" filter:"updated_at,string"`
	DB           *sql.DB  `json:"-"`
	Env          *env.Env `json:"-"`
}

func (user *User) Get(p *pagination.Paginator, f *filter.Filter) ([]*User, error) {
	var users []*User

	query := user.Env.QB.From(goqu.I("user")).Select(
		goqu.I("id"),
		goqu.I("email"),
		goqu.I("first_name"),
		goqu.I("last_name"),
		goqu.I("is_admin"),
		goqu.I("created_at"),
		goqu.I("updated_at")).Order(goqu.I("created_at").Desc()).Prepared(true)

	p.PK = "id"
	query = f.Filterize(query)
	query = p.Paginate(query)

	sqlstring, args, _ := query.ToSql()

	rows, err := user.Env.DB.Query(sqlstring, args...)
	defer rows.Close()
	for rows.Next() {
		u := new(User)
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt); err != nil {
			fmt.Printf("%+v\n", err)
		}
		users = append(users, u)
	}
	if err == nil {
		return users, nil
	} else if err == sql.ErrNoRows {
		return users, errors.New("there aren't any users")
	}
	return users, err
}

func (user *User) FindByEmail(email string) (*User, error) {
	err := user.DB.QueryRow("SELECT id, email as Email, password_hash FROM `user` WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err == nil {
		return user, nil
	} else if err == sql.ErrNoRows {
		return new(User), errors.New("there is no user with such email")
	}
	return new(User), err
}
func (user *User) Create() (int64, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	if err != nil {
		return 0, err
	}
	var isAdmin int
	if *user.IsAdmin == true {
		isAdmin = 1
	} else {
		isAdmin = 0
	}
	result, err := user.Env.DB.Exec("INSERT INTO user (`email`,`first_name`,`last_name`,`password_hash`,`is_admin`) VALUES (?,?,?,?,?) ", &user.Email, &user.FirstName, &user.LastName, bytes, isAdmin)

	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastID, nil
}
func (user *User) Update() error {

	var query string
	if user.Password != "" {
		query = "UPDATE user SET `first_name` = ?, `last_name` = ?, is_admin = ?, password_hash = ?  WHERE id=?"
	} else {
		query = "UPDATE user SET `first_name` = ?, `last_name` = ?, is_admin = ?  WHERE id=?"
	}
	sql, err := user.Env.DB.Prepare(query)
	if err != nil {
		return err
	}
	var isAdmin int
	if *user.IsAdmin == true {
		isAdmin = 1
	} else {
		isAdmin = 0
	}

	if user.Password != "" {
		bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		_, err = sql.Exec(&user.FirstName, &user.LastName, isAdmin, bytes, &user.ID)
	} else {
		_, err = sql.Exec(&user.FirstName, &user.LastName, isAdmin, &user.ID)
	}

	return err
}
func (user *User) GetByID(ID int64) (*User, error) {

	err := user.Env.DB.QueryRow("SELECT id, first_name, last_name, email, is_admin, created_at,updated_at FROM `user` where id = ? ", ID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) Delete() error {
	sql, err := user.Env.DB.Prepare("DELETE FROM user WHERE id=?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(&user.ID)
	return err
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}
