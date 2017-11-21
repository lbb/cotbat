package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
)

const (
	baseUrl    = "http://thecatapi.com"
	getPath    = "/api/images/get"
	configFile = "./config.json"
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

func printFancyMessageToLogBeansSingletonLoggerFacilityProviderBeans() {
	log.Println(`
  ______             __         ______               __
 /      \           /  |       /      \             /  |
/$$$$$$  |  ______  $$ |____  /$$$$$$  |  ______   _$$ |_
$$ |  $$/  /      \ $$      \ $$ |  $$/  /      \ / $$   |
$$ |      /$$$$$$  |$$$$$$$  |$$ |       $$$$$$  |$$$$$$/
$$ |   __ $$ |  $$ |$$ |  $$ |$$ |   __  /    $$ |  $$ | __
$$ \__/  |$$ \__$$ |$$ |__$$ |$$ \__/  |/$$$$$$$ |  $$ |/  |
$$    $$/ $$    $$/ $$    $$/ $$    $$/ $$    $$ |  $$  $$/
 $$$$$$/   $$$$$$/  $$$$$$$/   $$$$$$/   $$$$$$$/    $$$$/`)
	log.Println("Fancy app for fancy pep's!")
}

func parseValuesFromConfig() (port, imgType, imgSize string) {
	// Lookup port from os args
	portString := flag.Lookup(portKey).Value.String()
	absolutePath, err := filepath.Abs(configFile)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(absolutePath)
	if err != nil {
		log.Fatal(err)
	}
	// Close file after function exit
	defer f.Close()

	var i map[string]string
	// Parse file into i
	err = json.NewDecoder(f).Decode(&i)
	if err != nil {
		log.Fatal(err)
	}
	return portString, i[imgTypeKey], i[imgSizeKey]
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
	logString := flag.Lookup(logFileKey).Value.String()
	f, err := os.OpenFile(logString, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
}

func main() {
	// Add command line parameter
	flag.String(portKey, "80", "Specify port to serve on")
	flag.String(logFileKey, "log.log", "the output log file")
	flag.Parse()

	setupLogging()
	printFancyMessageToLogBeansSingletonLoggerFacilityProviderBeans()

	port, imgType, imgSize := parseValuesFromConfig()

	//go shutDownOnSignals()

	sv := buildServeFunction(baseUrl+getPath, imgType, imgSize)
	http.HandleFunc("/", sv)
	addr := ":" + port
	log.Printf("Start serving on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, http.DefaultServeMux))
}
