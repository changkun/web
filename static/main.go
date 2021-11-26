// Copyright 2021 Changkun Ou. All rights reserved.

package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var (
	app     *newrelic.Application
	wrapper = newrelic.WrapHandle
)

func init() {
	log.SetPrefix("static:")
	name := os.Getenv("NEWRELIC_NAME")
	lice := os.Getenv("NEWRELIC_LICENSE")

	if name == "" || lice == "" {
		// Don't use NewRelic is name or license is missing.
		wrapper = func(app *newrelic.Application, pattern string, handler http.Handler) (string, http.Handler) {
			return pattern, handler
		}
		log.Println("NewRelic is deactivated.")
		return
	}

	var err error
	app, err = newrelic.NewApplication(
		newrelic.ConfigAppName(name),
		newrelic.ConfigLicense(lice),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		log.Fatalf("Failed to created NewRelic application: %v", err)
	}

	log.Println("NewRelic is activated.")
}

func main() {
	log.SetPrefix("static:")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	l := logging()

	addr := os.Getenv("STATIC_ADDR")
	if addr == "" {
		log.Fatalf("missing address.")
	}
	folder := os.Getenv("STATIC_FOLDER")
	if !strings.HasPrefix(folder, "/www") {
		log.Fatal("failed to serve folder outside /www folder.")
	}

	http.Handle(wrapper(app, "/", l(http.FileServer(http.Dir(folder)))))

	log.Println("Listening on http://static:80...")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// logging wraps an http handler and returns a new handler that prints
// request log.
func logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				log.Println(readIP(r), r.Method, r.URL.Path, r.URL.RawQuery)
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// ReadIP implements a best effort approach to return the real client IP,
// it parses X-Real-IP and X-Forwarded-For in order to work properly with
// reverse-proxies such us: nginx or haproxy. Use X-Forwarded-For before
// X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
//
// The purpose of this function is to produce an identifier of visitor.
// It does not matter wheather it is an real IP or not. Depending on the
// configuration, the returned IP address might be an encrypted hash string.
//
// This implementation is derived from gin-gonic/gin.
func readIP(r *http.Request) (ip string) {
	ip = r.Header.Get("X-Forwarded-For")
	ip = strings.TrimSpace(strings.Split(ip, ",")[0])
	if ip == "" {
		ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	}
	if ip != "" {
		return ip
	}
	ip = r.Header.Get("X-Appengine-Remote-Addr")
	if ip != "" {
		return ip
	}
	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return "unknown" // use unknown to guarantee non empty string
	}
	return ip
}
