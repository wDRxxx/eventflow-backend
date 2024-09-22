package utils

import (
	"reflect"

	"github.com/pkg/errors"
)

// MapByStructTag takes tag value of each field
// as key for new map with appropriate field value
func MapByStructTag(tag string, data any) (m map[string]interface{}, err error) {
	defer func() {
		r := recover()
		if r != nil {
			m = nil

			recoverErr, ok := r.(error)
			if !ok {
				err = errors.New("error making map from struct")
				return
			}

			err = recoverErr
		}
	}()

	m = make(map[string]interface{})
	v := reflect.ValueOf(data)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		t := typeOfS.Field(i).Tag.Get(tag)
		if !v.Field(i).IsZero() && t != "" && t != "-" {
			m[t] = v.Field(i).Interface()
		}
	}

	return m, nil
}
