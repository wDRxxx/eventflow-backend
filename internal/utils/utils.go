package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"

	"github.com/google/uuid"
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
		if (!v.Field(i).IsZero() || typeOfS.Field(i).Type.String() == "bool") && t != "" && t != "-" {
			m[t] = v.Field(i).Interface()
		}
	}

	return m, nil
}

func SaveStaticImage(img *multipart.FileHeader, staticDir string) (string, error) {
	extension := filepath.Ext(img.Filename)
	imageName := uuid.New().String() + extension

	localFile, err := os.Create(staticDir + imageName)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	file, err := img.Open()
	_, err = io.Copy(localFile, file)
	defer file.Close()

	if err != nil {
		return "", err
	}

	return imageName, nil
}
