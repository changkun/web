// Copyright 2021 Changkun Ou. All rights reserved.

package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	addr := os.Getenv("STATIC_ADDR")
	if addr == "" {
		log.Fatalf("missing address.")
	}
	folder := os.Getenv("STATIC_FOLDER")
	if !strings.HasPrefix(folder, "/www") {
		log.Fatal("failed to serve folder outside /www folder.")
	}

	http.Handle("/", http.FileServer(http.Dir(folder)))
	log.Println("Listening on http://static:80...")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
