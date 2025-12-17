package main

import (
	"flag"
	"fmt"

	"github.com/Ericwyn/pancors"
)

const version = "0.2.0"

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

	port := pancors.GetListenPort(portFlag)
	pancors.RunPancorsServ(port)
}
