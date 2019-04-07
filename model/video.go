package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/arizanovj/courses/env"
	pagination "github.com/arizanovj/courses/libs"
	"github.com/arizanovj/courses/libs/filter"
	_ "github.com/go-sql-driver/mysql"
	goqu "gopkg.in/doug-martin/goqu.v4"
	_ "gopkg.in/doug-martin/goqu.v4/adapters/mysql"
)

type Video struct {
	ID          int64    `json:"id" filter:"id,number"`
	Name        string   `json:"name" filter:"name,string"`
	Description *string  `json:"description" filter:"description,string"`
	Cover       *string  `json:"cover" filter:"-"`
	Src         *string  `json:"src" filter:"-"`
	Offline     bool     `json:"offline" filter:"offline,string"`
	CourseID    int64    `json:"course_id" filter:"course,number"`
	CreatedAt   string   `json:"created_at"  filter:"created_at,date"`
	UpdatedAt   string   `json:"updated_at"  filter:"updated_at,date"`
	Env         *env.Env `json:"-"`
}

func (video *Video) Get(p *pagination.Paginator, f *filter.Filter) ([]*Video, error) {
	var videos []*Video

	query := video.Env.QB.From(goqu.I("video")).Order(goqu.I("created_at").Desc()).Prepared(true)

	p.PK = "id"
	query = f.Filterize(query)
	query = p.Paginate(query)

	sqlstring, args, _ := query.ToSql()

	rows, err := video.Env.DB.Query(sqlstring, args...)
	defer rows.Close()
	for rows.Next() {
		c := new(Video)
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Cover, &c.Src, &c.Offline, &c.CourseID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			fmt.Printf("%+v\n", err)
		}
		videos = append(videos, c)
	}
	if err == nil {
		return videos, nil
	} else if err == sql.ErrNoRows {
		return videos, errors.New("there aren't any videos")
	}
	return videos, err
}
func (video *Video) GetByID(ID int64) (*Video, error) {

	err := video.Env.DB.QueryRow("SELECT id, name, description, cover, src, course_id,offline, created_at,updated_at FROM video where id = ? ", ID).Scan(&video.ID, &video.Name, &video.Description, &video.Cover, &video.Src, &video.CourseID, &video.Offline, &video.CreatedAt, &video.UpdatedAt)
	if err != nil {
		return &Video{}, err
	}
	return video, nil

}

func (video *Video) Create() (int64, error) {

	result, err := video.Env.DB.Exec("INSERT INTO video (`name`,`description`,`course_id`,`offline`) VALUES (?,?,?,?) ", &video.Name, &video.Description, &video.CourseID, &video.Offline)

	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (video *Video) UpdateCover() error {
	sql, err := video.Env.DB.Prepare("UPDATE video SET cover=? WHERE id=?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(&video.Cover, &video.ID)

	return err
}

func (video *Video) UpdateSrc() error {
	sql, err := video.Env.DB.Prepare("UPDATE video SET src=? WHERE id=?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(&video.Src, &video.ID)

	return err
}

func (video *Video) Update() error {
	sql, err := video.Env.DB.Prepare("UPDATE video SET `name` = ?, `description` = ?,`offline` = ?  WHERE id=?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(&video.Name, &video.Description, &video.Offline, &video.ID)

	return err
}
func (video *Video) Delete() error {
	sql, err := video.Env.DB.Prepare("DELETE FROM video WHERE id=?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(&video.ID)
	return err
}
