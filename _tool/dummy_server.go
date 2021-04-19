package main

import (
	"fmt"
	"log"
	"net/http"
)

var count int

func main() {
	http.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "hello\n")
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error\n")
	})


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		count += 1

		if count%5 == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error\n")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "hello\n")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
