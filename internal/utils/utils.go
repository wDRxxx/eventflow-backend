package utils

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/pkg/errors"

	"github.com/wDRxxx/eventflow-backend/internal/models"
)

// MapByStructTags takes tag value of field as key for new map
// with appropriate value
func MapByStructTags(tag string, data any) (m map[string]interface{}, err error) {
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

func WriteJSON(data any, w http.ResponseWriter, status ...int) error {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	if len(status) > 0 {
		w.WriteHeader(status[0])
	}

	_, err = w.Write(resp)
	if err != nil {
		return err
	}

	return nil
}

// WriteJSONError writes json error
func WriteJSONError(err error, w http.ResponseWriter, status ...int) error {
	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(&models.DefaultResponse{
		Error:   true,
		Message: err.Error(),
	})

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusInternalServerError)
	if len(status) > 0 {
		w.WriteHeader(status[0])
	}

	_, err = w.Write(res)
	if err != nil {
		return err
	}

	return nil
}

// ReadJSON reads json to given data
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1024 * 1024 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&data)
	if err != nil {
		return err
	}

	return nil
}
