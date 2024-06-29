package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
    // handle paths
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
    // handle directories for static files
    handleDirectory("/js")

    // start the server
	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleDirectory(dirs... string) {
    for _, dir := range dirs {
        http.Handle(dir, http.StripPrefix(dir, http.FileServer(http.Dir("."+dir))))
    }
}
