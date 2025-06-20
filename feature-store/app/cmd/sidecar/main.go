package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("flashfeat sidecar start on :8080")
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	// TODO: vsock bridge to enclave
	log.Fatal(http.ListenAndServe(":8080", nil))
}
