package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/billhathaway/webcounter"
)

const (
	defaultPort = "8080"
)

func main() {
	port := flag.String("p", defaultPort, "listen port")
	flag.Parse()
	counter, err := webcounter.New()
	if err != nil {
		log.Fatal(err)
	}
	http.ListenAndServe(":"+*port, counter)
}
