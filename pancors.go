package pancors

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type corsTransport struct {
	referer     string
	origin      string
	credentials string
}

func (t corsTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	// Put in the Referer if specified
	if t.referer != "" {
		r.Header.Add("Referer", t.referer)
	}

	// Do the actual request
	res, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	res.Header.Set("Access-Control-Allow-Origin", t.origin)
	res.Header.Set("Access-Control-Allow-Credentials", t.credentials)

	// Ensure Mcp-Session-Id is exposed in CORS headers
	exposeHeaders := res.Header.Get("Access-Control-Expose-Headers")
	if exposeHeaders == "" {
		res.Header.Set("Access-Control-Expose-Headers", "Mcp-Session-Id")
	} else if !containsHeader(exposeHeaders, "Mcp-Session-Id") {
		res.Header.Set("Access-Control-Expose-Headers", exposeHeaders+", Mcp-Session-Id")
	}

	return res, nil
}

func handleProxy(w http.ResponseWriter, r *http.Request, origin string, credentials string) {
	// Handle OPTIONS requests directly
	if r.Method == "OPTIONS" {
		// Set CORS headers for OPTIONS preflight
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", credentials)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin, Access-Control-Request-Method, Access-Control-Request-Headers, Mcp-Session-Id")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check for the User-Agent header
	if r.Header.Get("User-Agent") == "" {
		http.Error(w, "Missing User-Agent header", http.StatusBadRequest)
		return
	}

	// Get the optional Referer header
	referer := r.URL.Query().Get("referer")
	if referer == "" {
		referer = r.Header.Get("Referer")
	}

	// Get the URL
	urlParam := r.URL.Query().Get("url")
	// Validate the URL
	urlParsed, err := url.Parse(urlParam)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	// Check if HTTP(S)
	if urlParsed.Scheme != "http" && urlParsed.Scheme != "https" {
		http.Error(w, "The URL scheme is neither HTTP nor HTTPS", http.StatusBadRequest)
		return
	}

	// Setup for the proxy
	proxy := httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL = urlParsed
			r.Host = urlParsed.Host
		},
		Transport: corsTransport{referer, origin, credentials},
	}

	// Execute the request
	proxy.ServeHTTP(w, r)
}

// HandleProxy is a handler which passes requests to the host and returns their
// responses with CORS headers
func HandleProxy(w http.ResponseWriter, r *http.Request) {
	handleProxy(w, r, "*", "true")
}

// HandleProxyWith is a handler which passes requests only from specified to the host
func HandleProxyWith(origin string, credentials string) func(http.ResponseWriter, *http.Request) {
	if !(credentials == "true" || credentials == "false") {
		log.Panicln("Access-Control-Allow-Credentials can only be 'true' or 'false'")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		handleProxy(w, r, origin, credentials)
	}
}

// containsHeader checks if a header name exists in a comma-separated list of headers
func containsHeader(headers, header string) bool {
	for _, h := range strings.Split(headers, ",") {
		if strings.TrimSpace(h) == header {
			return true
		}
	}
	return false
}
