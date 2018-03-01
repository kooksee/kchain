package validation

import (
	"reflect"
	"time"
	"gopkg.in/go-playground/validator.v8"
)

func BookableDate(
v *validator.Validate, topStruct reflect.Value,
currentStructOrField reflect.Value, field reflect.Value,
fieldType reflect.Type, fieldKind reflect.Kind, param string, ) bool {
	if date, ok := field.Interface().(time.Time); ok {
		today := time.Now()
		if today.Year() > date.Year() || today.YearDay() > date.YearDay() {
			return false
		}
	}
	return true
}