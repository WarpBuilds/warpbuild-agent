// cache_handler.go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// logs every cache-related path
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		// ← for a PoC just echo OK
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "cache request logged")

		// ↳ later: stream to MinIO, S3, or your real cache here
	})

	server := &http.Server{
		Addr:    "127.0.0.1:49161",
		Handler: mux,
	}

	log.Println("simple cache handler on :49161")
	log.Fatal(server.ListenAndServe())
}
