package model

import (
	"database/sql"

	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/go-ozzo/ozzo-validation"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DB       *sql.DB
}

func (l *Login) Login() (error, bool) {
	userModel := User{DB: l.DB}
	user, err := userModel.FindByEmail(l.Email)
	if err != nil {
		return err, false
	}
	err = user.ValidatePassword(l.Password)
	if err != nil {
		return err, false
	}

	return nil, true
}

func (l Login) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, validation.Required, validation.Length(5, 50), is.Email),
		validation.Field(&l.Password, validation.Required, validation.Length(8, 20)),
	)
}
