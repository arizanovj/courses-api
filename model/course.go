package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/arizanovj/courses/env"
	"github.com/arizanovj/courses/libs"
	"github.com/arizanovj/courses/libs/filter"
	_ "github.com/go-sql-driver/mysql"
	goqu "gopkg.in/doug-martin/goqu.v4"
	_ "gopkg.in/doug-martin/goqu.v4/adapters/mysql"
)

type Course struct {
	ID          int64    `json:"id" filter:"id,number"`
	Name        string   `json:"name" filter:"name,string"`
	Description *string  `json:"description" filter:"description,string"`
	Cover       *string  `json:"cover" filter:"-"`
	CreatedAt   string   `json:"created_at"  filter:"created_at,date"`
	UpdatedAt   string   `json:"updated_at"  filter:"updated_at,date"`
	Env         *env.Env `json:"-"`
}

func (course *Course) Get(p *pagination.Paginator, f *filter.Filter) ([]*Course, error) {
	var courses []*Course

	query := course.Env.QB.From(goqu.I("course")).Order(goqu.I("created_at").Desc()).Prepared(true)

	p.PK = "id"
	query = f.Filterize(query)
	query = p.Paginate(query)

	sqlstring, args, _ := query.ToSql()

	rows, err := course.Env.DB.Query(sqlstring, args...)
	defer rows.Close()
	for rows.Next() {
		c := new(Course)
		if err := rows.Scan(&c.ID, &c.Name, &c.Cover, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			fmt.Printf("%+v\n", err)
		}
		courses = append(courses, c)
	}
	if err == nil {
		return courses, nil
	} else if err == sql.ErrNoRows {
		return courses, errors.New("there aren't any courses")
	}
	return courses, err
}
func (course *Course) GetByID(ID int64) (*Course, error) {

	err := course.Env.DB.QueryRow("SELECT id, name, description, cover, created_at,updated_at FROM course where id = ? ", ID).Scan(&course.ID, &course.Name, &course.Description, &course.Cover, &course.CreatedAt, &course.UpdatedAt)
	if err != nil {
		return &Course{}, err
	}
	return course, nil

}

func (course *Course) Create() (int64, error) {

	result, err := course.Env.DB.Exec("INSERT INTO course (`name`,`description`) VALUES (?,? ) ", &course.Name, &course.Description)

	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (course *Course) UpdateCover() error {
	sql, err := course.Env.DB.Prepare("UPDATE course SET cover=? WHERE id=?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(&course.Cover, &course.ID)

	return err
}

func (course *Course) Update() error {
	sql, err := course.Env.DB.Prepare("UPDATE course SET `name` = ?, `description` = ?  WHERE id=?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(&course.Name, &course.Description, &course.ID)

	return err
}
func (course *Course) Delete() error {
	sql, err := course.Env.DB.Prepare("DELETE FROM course WHERE id=?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(&course.ID)
	return err
}
