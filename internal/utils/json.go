package utils

import (
	"encoding/json"
	"net/http"

	"github.com/wDRxxx/eventflow-backend/internal/models"
)

// WriteJSON writes json of data to w
func WriteJSON(data any, w http.ResponseWriter, status ...int) error {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(status) == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(status[0])
	}

	_, err = w.Write(resp)
	if err != nil {
		return err
	}

	return nil
}

// WriteJSONError writes json error to w
func WriteJSONError(err error, w http.ResponseWriter, status ...int) error {
	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(&models.DefaultResponse{
		Error:   true,
		Message: err.Error(),
	})

	if err != nil {
		return err
	}

	if len(status) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
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
