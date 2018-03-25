package weque

import (
	"log"
	"net/http"
)

func SendError(w http.ResponseWriter, status int, msg string) {
	log.Print(msg)
	w.WriteHeader(status)
	w.Write([]byte(msg))
}
