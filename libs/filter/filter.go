package filter

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/arizanovj/courses/env"
	"github.com/go-ozzo/ozzo-validation"
	goqu "gopkg.in/doug-martin/goqu.v4"
)

const (
	Contains             = "cnt"
	ContainsEnd          = "cnte"
	ContainsStart        = "cnts"
	EqualTo              = "eq"
	LessThan             = "lt"
	LessThanOrEqualTo    = "lte"
	GreaterThan          = "gt"
	GreaterThanOrEqualTo = "gte"
)

var Types = [...]string{
	"number",
	"date",
	"string",
}

var NumberFilters = []string{
	Contains,
	ContainsEnd,
	ContainsStart,
	EqualTo,
	LessThan,
	LessThanOrEqualTo,
	GreaterThan,
	GreaterThanOrEqualTo,
}

var StringFilters = []string{
	Contains,
	ContainsEnd,
	ContainsStart,
	EqualTo,
}
var FromDateFilters = []string{
	GreaterThan,
	GreaterThanOrEqualTo,
}

var ToDateFilters = []string{
	LessThan,
	LessThanOrEqualTo,
}

var tagName = "filter"

type Filter struct {
	filterParams map[string][]string
	Model        interface{}
	Errors       []error
	Env          *env.Env
}

func (f *Filter) Filterize(query *goqu.Dataset) *goqu.Dataset {

	for key, value := range f.filterParams {
		field := strings.TrimSuffix(strings.TrimPrefix(key, "filter["), "]")
		k := strings.Split(field, "|")
		v := strings.Split(value[0], "|")

		if len(k) == 1 && f.validateField(k[0], v) {
			query = query.Where(f.getQuery(string(k[0]), string(v[0]), string(v[1])))

		} else if f.validateDateField(k, v) {
			query = query.Where(f.getQuery(string(k[0]), string(v[0]), string(v[1])))
		}
	}

	return query
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (f *Filter) validateDateField(key []string, value []string) bool {

	fieldType := f.getFieldTypeFromTag(key[0])

	if fieldType == "date" {
		if key[1] == "from" && stringInSlice(value[0], FromDateFilters) {
			return true
		}

		if key[1] == "to" && stringInSlice(value[0], ToDateFilters) {
			return true
		}

	}
	f.Errors = append(f.Errors, errors.New("field is not date type"+key[0]))
	return false
}
func (f *Filter) validateField(key string, value []string) bool {
	fieldType := f.getFieldTypeFromTag(key)

	if fieldType == "number" && stringInSlice(value[0], NumberFilters) {
		return true
	}

	if fieldType == "string" && stringInSlice(value[0], StringFilters) {
		return true
	}

	f.Errors = append(f.Errors, errors.New("invalid data per "+key))
	return false
}

func (f *Filter) getQuery(field string, filter string, value string) goqu.Expression {

	switch filter {
	case "gt":
		return goqu.I(field).Gt(value)
	case "gte":
		return goqu.I(field).Gte(value)
	case "lt":
		return goqu.I(field).Lt(value)
	case "lte":
		return goqu.I(field).Lte(value)
	case "eq":
		return goqu.I(field).Eq(value)
	case "cnt":
		return goqu.I(field).Like("%" + value + "%")
	case "cnts":
		return goqu.I(field).Like(value + "%")
	case "cnte":
		return goqu.I(field).Like("%" + value)
	}
	return goqu.I("1").Eq("1")
}

func (f *Filter) getFieldTypeFromTag(field string) string {
	t := reflect.TypeOf(f.Model)

	for i := 0; i < t.NumField(); i++ {

		f := t.Field(i)
		tag := f.Tag.Get(tagName)

		if tag == "" || tag == "-" {
			continue
		}
		filterTag := strings.Split(tag, ",")

		if string(filterTag[0]) == field {
			return string(filterTag[1])
		}

	}
	return ""
}

func (f *Filter) SetFilterParams(requestParams map[string][]string) {
	for key, value := range requestParams {

		if len(value) > 1 || f.isFilterKey(key) != nil || f.isFilterValue(value[0]) != nil {
			delete(requestParams, key)
		}
	}
	f.filterParams = requestParams
}

func (f *Filter) isFilterKey(key string) error {
	return validation.Validate(key,
		validation.Match(regexp.MustCompile("^filter\\[([a-z0-9_])+(\\|from|\\|to)?]")),
	)
}

func (f *Filter) isFilterValue(value string) error {

	return validation.Validate(value,
		validation.Match(regexp.MustCompile("^("+strings.Join(f.GetFilters(), "|")+"){1}\\|[A-Za-z0-9_]+")),
	)
}

func (f *Filter) GetFilters() []string {
	s := []string{
		Contains,
		ContainsEnd,
		ContainsStart,
		EqualTo,
		LessThan,
		LessThanOrEqualTo,
		GreaterThan,
		GreaterThanOrEqualTo,
	}
	return s
}
