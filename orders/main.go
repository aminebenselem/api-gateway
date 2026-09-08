package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	instance := os.Getenv("INSTANCE_NAME")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(
			w,
			"Hello from %s: %s %s",
			instance,
			r.Method,
			r.URL.Path,
		)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	fmt.Printf("Orders service %s listening on :8080\n", instance)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}