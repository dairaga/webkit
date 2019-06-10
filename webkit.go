package webkit

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var _router = mux.NewRouter()
var _subrouters = make(map[string]*mux.Router)

// Handle set http.Handler for url.
func Handle(url string, f interface{}, params ...string) *mux.Route {
	for k, r := range _subrouters {
		if url == k || strings.HasPrefix(url, k+"/") {
			return r.Handle(strings.Replace(url, k, "", -1), _mkHandleFunc(f, params...))
		}
	}
	return _router.Handle(url, _mkHandleFunc(f, params...))
}

// Use returns a subrouter.
func Use(prefix string, filters ...Filter) *mux.Router {
	var r *mux.Router

	if prefix != "/" {
		var ok bool
		r, ok = _subrouters[prefix]
		if !ok {
			r = _router.PathPrefix(prefix).Subrouter()
		}

		_subrouters[prefix] = r
	} else {
		r = _router
	}

	for _, f := range filters {
		r.Use(_mkMiddlewareFunc(f))
	}
	return r
}

// Router returns inside router.
func Router() *mux.Router {
	return _router
}

// Start run as http server.
func Start(host ...string) error {
	if len(host) <= 0 || host[0] == "" {
		return http.ListenAndServe(":80", _router)
	}

	return http.ListenAndServe(host[0], _router)
}

// StartSecure run as https.
func StartSecure(cert, key string, host ...string) error {
	if len(host) <= 0 || host[0] == "" {
		return http.ListenAndServeTLS(":443", cert, key, _router)
	}

	return http.ListenAndServeTLS(host[0], cert, key, _router)
}
