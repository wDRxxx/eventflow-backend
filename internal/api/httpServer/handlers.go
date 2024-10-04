package httpServer

import (
	"net/http"

	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

func saveMultipartImages(r *http.Request, formField string, staticDir string) ([]string, error) {
	reqImages := r.MultipartForm.File[formField]
	var images []string

	for _, img := range reqImages {
		filename, err := utils.SaveStaticImage(img, staticDir)
		if err != nil {
			return nil, err
		}

		images = append(images, filename)
	}

	return images, nil
}
