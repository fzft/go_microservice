package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Hello world")
		d, err := ioutil.ReadAll(r.Body)
		if err !=nil {
			http.Error(rw, "Ooops", http.StatusBadRequest)
			return
		}
		log.Printf("Data %s\n", d)
		fmt.Fprintf(rw, "Hello %s\n", d)

	})

	http.HandleFunc("/goodbye", func(w http.ResponseWriter, _ *http.Request) {
		log.Println("Goodbye world")
	})

	http.ListenAndServe(":9090", nil)
}
