package api

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/evilsocket/islazy/log"
)

// DNSResolver interface for reverse DNS lookups, allows mocking in tests
type DNSResolver interface {
	LookupAddr(ctx context.Context, addr string) ([]string, error)
}

// defaultResolver uses net.DefaultResolver for DNS lookups
type defaultResolver struct{}

func (r *defaultResolver) LookupAddr(ctx context.Context, addr string) ([]string, error) {
	return net.DefaultResolver.LookupAddr(ctx, addr)
}

// DefaultDNSResolver is the default resolver instance
var DefaultDNSResolver DNSResolver = &defaultResolver{}

var (
	ErrEmpty           = errors.New("")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrMessageNotFound = errors.New("record not found")
)

func strParam(r *http.Request, name string, def string) string {
	param := r.URL.Query().Get(name)
	if param != "" {
		return param
	}
	return def
}

func floatParam(r *http.Request, name string, def float32) float32 {
	param := r.URL.Query().Get(name)
	if param != "" {
		if v, err := strconv.ParseFloat(param, 32); err != nil {
			log.Warning("can't parse %s parameter from value '%s'", name, param)
		} else {
			return float32(v)
		}
	}
	return def
}

func boolParam(r *http.Request, name string, def bool) bool {
	param := r.URL.Query().Get(name)
	if param != "" {
		if v, err := strconv.ParseBool(param); err != nil {
			log.Warning("can't parse %s parameter from value '%s'", name, param)
		} else {
			return v
		}
	}
	return def
}

func intParam(r *http.Request, name string, def int) int {
	param := r.URL.Query().Get(name)
	if param != "" {
		if v, err := strconv.ParseInt(param, 10, 32); err != nil {
			log.Warning("can't parse %s parameter from value '%s'", name, param)
		} else {
			return int(v)
		}
	}
	return def
}

func clientIP(r *http.Request) string {
	address := strings.Split(r.RemoteAddr, ":")[0]
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		address = forwardedFor
	}
	// https://support.cloudflare.com/hc/en-us/articles/206776727-What-is-True-Client-IP-
	if trueClient := r.Header.Get("True-Client-IP"); trueClient != "" {
		address = trueClient
	}
	// handle multiple IPs case
	return strings.Trim(strings.Split(address, ",")[0], " ")
}

func reqToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if parts := strings.Split(bearerToken, " "); len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func pageNum(r *http.Request) (int, error) {
	pageParam := r.URL.Query().Get("p")
	if pageParam == "" {
		pageParam = "1"
	}
	return strconv.Atoi(pageParam)
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err = w.Write(js); err != nil {
		log.Error("error sending response: %v", err)
	} else {
		// log.Debug("sent %d bytes of json response", sent)
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}

// ResolveHostname performs a reverse DNS lookup for the given IP address
// and returns the hostname. If the lookup fails or times out, it returns
// an empty string. The resolver parameter allows for dependency injection
// in tests; if nil, the default resolver is used.
func ResolveHostname(ipAddress string, resolver DNSResolver) string {
	if resolver == nil {
		resolver = DefaultDNSResolver
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	names, err := resolver.LookupAddr(ctx, ipAddress)
	if err != nil {
		log.Debug("reverse DNS lookup failed for %s: %v", ipAddress, err)
		return ""
	}

	if len(names) > 0 {
		hostname := strings.TrimSuffix(names[0], ".")
		return hostname
	}

	return ""
}
