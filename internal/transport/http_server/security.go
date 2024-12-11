package http_server

import (
	"net/http"
	"strings"
)

// Get client IP-address
func GetClientIP(r *http.Request) string {
	//TODO: strip port from return value
	// check X-Forwarded-For header if we sit behing proxy/balancer
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// Get first IP if there's more than one
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// Use RemoteAddr if there's no header
	ip = r.RemoteAddr
	return ip
}

func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	return err == nil && cookie.Value == "authenticated"
}
