package handlers

import (
	"net/http"
)

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("This is an example endpoint!"))
	if err != nil {
		return
	}
}
