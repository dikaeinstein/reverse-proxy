package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

type Proxy interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type ReverseProxyHandler struct {
	target *url.URL
	proxy  Proxy
}

func (rph ReverseProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	r.URL.Host = rph.target.Host
	r.URL.Scheme = rph.target.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = rph.target.Host

	rph.proxy.ServeHTTP(w, r)
}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	port, _ := os.LookupEnv("PORT")
	u := flag.String("host", "", "backend URL host to proxy request to")
	flag.Parse()

	if *u == "" {
		*u, _ = os.LookupEnv("TARGET_HOST")
	}

	fmt.Println("Target host: ", *u)
	target, _ := url.Parse(*u)
	reverseProxyHandler := ReverseProxyHandler{
		target: target,
		proxy:  httputil.NewSingleHostReverseProxy(target),
	}
	http.Handle("/", reverseProxyHandler)

	log.Printf("Proxy server listening on port: %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
