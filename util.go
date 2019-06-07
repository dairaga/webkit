package webkit

import (
	"net/http"
	"strings"
)

// ClientIP returns client ip.
func ClientIP(r *http.Request) string {
	tmp := r.Header.Get("HTTP_X_FORWARDED_FOR")
	if tmp != "" {
		ips := strings.Split(tmp, ",")
		ipsLen := len(ips)
		if len(ips) > 0 {
			return strings.TrimSpace(ips[ipsLen-1])
		}
	}

	tmp = r.Header.Get("HTTP_X_REAL_IP")
	if tmp != "" {
		return tmp
	}
	tmp = r.Header.Get("HTTP_CLIENT_IP")
	if tmp != "" {
		return tmp
	}

	tmp = r.Header.Get("REMOTE_ADDR")
	if tmp != "" {
		return tmp
	}
	return r.RemoteAddr
}
