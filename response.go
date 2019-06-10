package webkit

import (
	"net/http"
)

// Result ...
type Result interface {
	Code() int
	Data() interface{}
	Templates() []string
}

type _resultImpl struct {
	code  int
	data  interface{}
	tmpls []string
}

func (r *_resultImpl) Code() int {
	return r.code
}

func (r *_resultImpl) Data() interface{} {
	return r.data
}

func (r *_resultImpl) Templates() []string {
	return r.tmpls
}

// Custom returns a custom result.
func Custom(code int, data interface{}, tmpls ...string) Result {
	return &_resultImpl{
		code:  code,
		data:  data,
		tmpls: tmpls,
	}
}

// OK (200)
func OK(data interface{}, tmpls ...string) Result {
	return Custom(http.StatusOK, data, tmpls...)
}

// Accepted (202)
func Accepted(data interface{}, tmpls ...string) Result {
	return Custom(http.StatusAccepted, data, tmpls...)
}

// NoContent (204)
func NoContent() Result {
	return Custom(http.StatusNoContent, nil)
}

// Redirect redriect with default 302...
func Redirect(w http.ResponseWriter, url string, code ...int) Result {
	status := http.StatusFound

	if len(code) > 0 {
		status = code[0]
	}

	w.Header().Add("Location", url)
	w.WriteHeader(status)
	return nil
}

// BadRequest (400)
func BadRequest(data interface{}, tmpls ...string) Result {
	return Custom(http.StatusBadRequest, data, tmpls...)
}

// Unauthorized (401)
func Unauthorized(data interface{}, tmpls ...string) Result {
	return Custom(http.StatusUnauthorized, data, tmpls...)
}

// Forbidden (403)
func Forbidden(data interface{}, tmpls ...string) Result {
	return Custom(http.StatusForbidden, data, tmpls...)
}

// NotFound (404)
func NotFound(data interface{}, tmpls ...string) Result {
	return Custom(http.StatusNotFound, data, tmpls...)
}

// Conflict (409)
func Conflict(data interface{}, tmpls ...string) Result {
	return Custom(http.StatusConflict, data, tmpls...)
}

// Fault (500)
func Fault(data interface{}, tmpls ...string) Result {
	return Custom(http.StatusInternalServerError, data, tmpls...)
}
