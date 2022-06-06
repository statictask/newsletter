package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/statictask/newsletter/pkg/subscription"
)

func main() {
	subsc := subscription.New()
	subsc.Email = "test@test.tt"

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, fmt.Sprintf("Hello, %s!\n", subsc.Email))
	}

	http.HandleFunc("/subscriptions", helloHandler)
	log.Println("Listing for requests at http://localhost:8000/subscriptions")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
