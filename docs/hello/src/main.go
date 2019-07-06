package main

import (
	"github.com/sipin/gorazor/docs/hello/src/tpl"

	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, tpl.Index())
	})

	http.ListenAndServe(":8080", nil)
}
