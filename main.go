package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
)

const (
	baseUrl    = "http://thecatapi.com"
	getPath    = "/api/images/get"
)

const (
	portKey    = "port"
	logFileKey = "log"

	// Can be:
	// jpg, png, gif
	imgTypeKey = "type"

	// Can be:
	// small, med, full
	imgSizeKey = "size"
)

func shutDownOnSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch) //,syscall.SIGTERM, syscall.SIGKILL)
	<-ch
	os.Exit(0)
}

func parseValuesFromConfig() (port, imgType, imgSize string) {
	// Lookup port from os args
	portString := flag.Lookup(portKey).Value.String()
	return portString, os.Getenv("COT_TYPE"), os.Getenv("COT_SIZE")
}

func buildServeFunction(urlPath, imgType, imgSize string) http.HandlerFunc {
	u, err := url.Parse(urlPath)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set(imgTypeKey, imgType)
	q.Set(imgSizeKey, imgSize)
	u.RawQuery = q.Encode()
	fullUrlPath := u.String()
	log.Printf("Rendered cat-pic URL: %s\n", fullUrlPath)
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got user connection: %s\n", r.Header.Get("User-Agent"))
		http.Redirect(w, r, fullUrlPath, http.StatusPermanentRedirect)
	}
}

func setupLogging() {
}

func main() {
	// Add command line parameter
	flag.String(portKey, "80", "Specify port to serve on")
	flag.String(logFileKey, "log.log", "the output log file")
	flag.Parse()

	setupLogging()

	port, imgType, imgSize := parseValuesFromConfig()

	//go shutDownOnSignals()

	sv := buildServeFunction(baseUrl+getPath, imgType, imgSize)
	http.HandleFunc("/", sv)
	addr := ":" + port
	log.Printf("Start serving on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, http.DefaultServeMux))
}
