package main

import (
	"fmt"
	"log"
	"net/http"
	// "log"
	// "net/http"
)

func main() {
	fmt.Println("Hello Fucking World")

	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte("<h1>Hello Fucking World<h1>"))
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/login", fileServer)

	log.Fatal(srv.ListenAndServe())

}
