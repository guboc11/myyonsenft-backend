package main

import (
	"fmt"
	"log"
	"net/http"
)

func mint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mint success"))

}

func main() {
	fmt.Println("hello Planetarium")

	http.HandleFunc("/mint", mint)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
