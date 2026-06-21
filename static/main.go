// Copyright 2021 Changkun Ou. All rights reserved.

package main

import (
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	log.SetPrefix("static: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	l := logging()

	addr := os.Getenv("STATIC_ADDR")
	if addr == "" {
		log.Fatalf("missing address.")
	}

	vanityImport := os.Getenv("STATIC_VANITY_IMPORT")
	if vanityImport != "" {
		http.Handle("/", l(vanityHandler(
			vanityImport,
			os.Getenv("STATIC_VANITY_SOURCE"),
			os.Getenv("STATIC_VANITY_REDIRECT"),
		)))
	} else {
		folder := os.Getenv("STATIC_FOLDER")
		if !strings.HasPrefix(folder, "/www") {
			log.Fatal("failed to serve folder outside /www folder.")
		}

		http.Handle("/", l(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow CORS for all /demos requests
			if strings.Contains(r.URL.Path, "/demos") {
				w.Header().Set("Access-Control-Allow-Headers", "*")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			}

			http.FileServer(http.Dir(folder)).ServeHTTP(w, r)
		})))
	}

	log.Println("Listening on http://static:80...")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func vanityHandler(importMeta, sourceMeta, redirectURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("go-get") == "1" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, "<!doctype html>\n<html><head>\n")
			fmt.Fprintf(w, "<meta name=\"go-import\" content=\"%s\">\n", html.EscapeString(importMeta))
			if sourceMeta != "" {
				fmt.Fprintf(w, "<meta name=\"go-source\" content=\"%s\">\n", html.EscapeString(sourceMeta))
			}
			fmt.Fprintf(w, "</head><body></body></html>\n")
			return
		}

		if redirectURL != "" {
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
			return
		}
		http.NotFound(w, r)
	})
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
