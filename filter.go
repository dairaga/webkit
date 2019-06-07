package webkit

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Filter a function for mux.MiddlewareFunc.
type Filter func(http.ResponseWriter, *http.Request) (*http.Request, bool)

func _mkMiddlewareFunc(f Filter) mux.MiddlewareFunc {

	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if newR, ok := f(w, r); ok {
				next.ServeHTTP(w, newR)
			}
		})
	})
}
