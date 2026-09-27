// Copyright 2021 Changkun Ou. All rights reserved.

package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
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

		var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow CORS for all /demos requests
			if strings.Contains(r.URL.Path, "/demos") {
				w.Header().Set("Access-Control-Allow-Headers", "*")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			}

			http.FileServer(http.Dir(folder)).ServeHTTP(w, r)
		})
		// STATIC_URLSTAT names a folder whose files are counted, as the
		// URL urlstat records them under, e.g. https://changkun.de/research/.
		if counted := os.Getenv("STATIC_URLSTAT"); counted != "" {
			u, err := url.Parse(counted)
			if err != nil || u.Scheme == "" || u.Host == "" || !strings.HasSuffix(u.Path, "/") {
				log.Fatalf("STATIC_URLSTAT must be the URL of a folder, got %q", counted)
			}
			h = countedHandler(u.Path, u.Scheme+"://"+u.Host, func(page, ua string) {
				go reportUrlstat(page, ua)
			}, h)
		}
		http.Handle("/", l(h))
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

// countedHandler reports every request below prefix to urlstat, as the
// page at origin+path, except the folder's own page, which a browser
// reports through urlstat's client script. What it counts is what a
// script cannot: the files a page links to, such as the PDFs of
// changkun.de/research.
func countedHandler(prefix, origin string, report func(page, ua string), next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if strings.HasPrefix(p, prefix) && p != prefix && p != prefix+"index.html" {
			report(origin+p, r.UserAgent())
		}
		next.ServeHTTP(w, r)
	})
}

// urlstatEndpoint receives the reports, the way urlstat's client script
// sends them from a browser.
var (
	urlstatEndpoint = "https://www.changkun.de/urlstat?report=page+site"
	urlstatClient   = &http.Client{Timeout: 10 * time.Second}
)

func reportUrlstat(page, ua string) {
	req, err := http.NewRequest(http.MethodGet, urlstatEndpoint, nil)
	if err != nil {
		log.Printf("urlstat: cannot report %s: %v", page, err)
		return
	}
	req.Header.Set("urlstat-url", page)
	req.Header.Set("urlstat-ua", ua)
	resp, err := urlstatClient.Do(req)
	if err != nil {
		log.Printf("urlstat: cannot report %s: %v", page, err)
		return
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("urlstat: cannot report %s: %s", page, resp.Status)
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
