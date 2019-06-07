package webkit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/dairaga/config"
	"github.com/dairaga/log"
	"github.com/gorilla/securecookie"
)

var _secureC *securecookie.SecureCookie
var _cookieDomain string

func init() {
	_cookieDomain = config.GetString("cookie.domain", "")
	hashK := config.GetString("cookie.hash", "")
	blockK := config.GetString("cookie.block", "")
	if hashK != "" && blockK != "" {
		_secureC = securecookie.New([]byte(hashK), []byte(blockK))
	}

	if _secureC != nil {
		log.Debug("init secure cookie")
	} else {
		log.Debug("init cookie")
	}
}

// UseSecureCookie ...
func UseSecureCookie(hashKey, blockKey string) {
	_secureC = securecookie.New([]byte(hashKey), []byte(hashKey))
}

// CookieDomain ...
func CookieDomain(domain string) {
	_cookieDomain = domain
}

func _encode(name string, value interface{}) (string, error) {
	if _secureC != nil {
		return _secureC.Encode(name, value)
	}

	var valuebytes []byte
	switch v := value.(type) {
	case string:
		valuebytes = []byte(v)
	case fmt.Stringer:
		valuebytes = []byte(v.String())
	default:
		var err error
		valuebytes, err = json.Marshal(value)
		if err != nil {
			return "", err
		}
	}

	return base64.StdEncoding.EncodeToString(valuebytes), nil
}

func _decode(name, value string, data interface{}) error {
	if _secureC != nil {
		return _secureC.Decode(name, value, data)
	}

	databytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return err
	}

	switch v := data.(type) {
	case *string:
		*v = string(databytes)
		return nil
	default:
		return json.Unmarshal(databytes, data)
	}
}

func _cookie(name, val, domain, path string, age int, secure, httpOnly bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    val,
		Domain:   domain,
		Path:     path,
		MaxAge:   age,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
}

// NewCookie ...
func NewCookie(name string, value interface{}, domain, path string, age int, secure, httpOnly bool) (*http.Cookie, error) {
	if value == nil {
		return _cookie(name, "", domain, path, age, secure, httpOnly), nil
	}

	val, err := _encode(name, value)
	if err != nil {
		return nil, err
	}

	return _cookie(name, val, domain, path, age, secure, httpOnly), nil
}

// DeleteCookie ...
func DeleteCookie(w http.ResponseWriter, name string) {
	SetCookie(w, name, nil, -1)
}

// SetCookie ...
func SetCookie(w http.ResponseWriter, name string, value interface{}, maxAge int) error {
	c, err := NewCookie(name, value, _cookieDomain, "/", maxAge, false, false)
	if err != nil {
		return err
	}

	http.SetCookie(w, c)
	return nil
}

// Cookie ...
func Cookie(r *http.Request, name string, data interface{}) error {
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return fmt.Errorf("data must be pointer of something")
	}

	c, err := r.Cookie(name)
	if err != nil {
		return err
	}

	return _decode(name, c.Value, data)
}
