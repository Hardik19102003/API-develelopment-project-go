package main

import (
	"fmt"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func main() {
	http.HandleFunc("/yup", homePage)
	http.ListenAndServe("localhost:8080", nil)
}
