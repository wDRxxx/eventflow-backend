package httpServer

import (
	"log"
	"net/http"
)

func (s *server) home(w http.ResponseWriter, r *http.Request) {
	log.Println("zxczxc")
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
