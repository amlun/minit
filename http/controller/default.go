package controller

import "net/http"

// DefaultHandler handles / request
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!\n"))
}
