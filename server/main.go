package main

import (
	"flag"
	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
	"github.com/wdongyu/builder-manager/server/handler"
	"net/http"
	"os"
	"strconv"
)

const (
	DefaultListenPort = 8088
)

func main() {
	var listenPort int
	flag.IntVar(&listenPort, "port", DefaultListenPort, `port this server listen on`)
	flag.Parse()

	h, err := handler.NewAPIHandler()
	if err != nil {
		log.Fatalf("Failed to create route handler: %v", err)
	}

	http.Handle("/", handler.NewRouter(h))
	log.Infof("Server listens on: %v", listenPort)
	if err = http.ListenAndServe(":"+strconv.Itoa(listenPort), handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}
