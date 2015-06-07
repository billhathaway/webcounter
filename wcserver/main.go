package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/billhathaway/webcounter"
)

const (
	defaultPort = "8080"
)

func main() {
	port := flag.String("p", defaultPort, "listen port")
	pprofPort := flag.String("pprof", "", "listen port for profiling")
	flag.Parse()
	counter, err := webcounter.New()
	if err != nil {
		log.Fatal(err)
	}
	if *pprofPort != "" {
		go func() {
			log.Fatal(http.ListenAndServe(":"+*pprofPort, nil))
		}()
	}
	log.Fatal(http.ListenAndServe(":"+*port, counter))
}
