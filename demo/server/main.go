package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Listening on :8080")
	if err := http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`))); err != nil {
		panic(err.Error())
	}
}
