package pagination

import (
	"fmt"

	"github.com/arizanovj/courses/env"
	validation "github.com/go-ozzo/ozzo-validation"
	_ "github.com/go-sql-driver/mysql"
	goqu "gopkg.in/doug-martin/goqu.v4"
	_ "gopkg.in/doug-martin/goqu.v4/adapters/mysql"
)

type Paginator struct {
	LastId     *uint64 `schema:"lastId" json:"lastId"`
	Direction  string  `schema:"direction" json:"direction"`
	NumOfItems int64   `schema:"numOfItems" json:"numOfItems"`
	PK         string
	Env        *env.Env `schema:"-" json:"-"`
}

func (p *Paginator) Paginate(query *goqu.Dataset) *goqu.Dataset {
	order := query.GetClauses().Order

	if order != nil {
		query.ClearOrder()
	}
	query = query.Order(p.GetOrder())
	if order != nil {
		query.GetClauses().Order.Append(order)
	}
	if p.LastId != nil {
		fmt.Printf("%+v\n", order)
		query = query.Where(p.GetWhere())
	}

	query = query.Limit(uint(p.NumOfItems))

	if p.Direction == "down" {
		query = p.Invert(query)
	}

	return query
}

func (p *Paginator) GetOrder() goqu.OrderedExpression {
	if p.Direction == "up" {
		return goqu.I(p.PK).Asc()
	} else {
		return goqu.I(p.PK).Desc()
	}
}

func (p *Paginator) GetWhere() goqu.Expression {
	if p.Direction == "up" {
		return goqu.I(p.PK).Gt(p.LastId)
	} else {
		return goqu.I(p.PK).Lt(p.LastId)
	}
}

func (p *Paginator) Invert(query *goqu.Dataset) *goqu.Dataset {
	return p.Env.QB.From(query).Order(goqu.I(p.PK).Asc())
}

func (p Paginator) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.LastId),
		validation.Field(&p.NumOfItems, validation.Required, validation.Min(5), validation.Max(100)),
		validation.Field(&p.Direction, validation.Required, validation.In("up", "down")),
	)
}
