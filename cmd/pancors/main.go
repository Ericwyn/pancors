package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/michaljanocko/pancors"
)

func getAllowOrigin() string {
	if origin, ok := os.LookupEnv("ALLOW_ORIGIN"); ok {
		return origin
	}
	return "*"
}

func getAllowCredentials() string {
	if credentials, ok := os.LookupEnv("ALLOW_CREDENTIALS"); ok {
		return credentials
	}
	return "true"
}

func getListenPort(portFlag string) string {
	// Command line flag takes precedence over environment variable
	if portFlag != "" {
		// Ensure port starts with colon
		if portFlag[0] != ':' {
			return ":" + portFlag
		}
		return portFlag
	}

	if port, ok := os.LookupEnv("PORT"); ok {
		return ":" + port
	}
	return ":8080"
}

const version = "0.1.0"

func main() {
	// Define command line flags
	var showVersion bool
	var portFlag string

	flag.BoolVar(&showVersion, "v", false, "Show version information")
	flag.StringVar(&portFlag, "port", "", "Port to listen on (overrides PORT environment variable)")
	flag.Parse()

	// Handle version flag
	if showVersion {
		fmt.Printf("PanCORS version %s\n", version)
		return
	}

	http.HandleFunc("/", pancors.HandleProxyWith(getAllowOrigin(), getAllowCredentials()))

	port := getListenPort(portFlag)
	log.Printf("PanCORS started listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
