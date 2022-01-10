package app

import (
	"fmt"
	"net/http"
)

func Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("in the middleware", r.URL)
		h(w, r)
	}
}
