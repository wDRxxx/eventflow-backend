package api

import "net/http"

type HTTPServer interface {
	Handler() http.Handler
}
