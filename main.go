package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func handlerHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = ""
	r.Header.Del("Proxy-Connection")

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	CopyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func handlerHTTPS(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)

}
func transfer(writeCloser io.WriteCloser, readCloser io.ReadCloser) {
	defer writeCloser.Close()
	defer readCloser.Close()
	_, _ = io.Copy(writeCloser, readCloser)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INIT")
	fmt.Println(r.Method)
	fmt.Println(r.Proto)
	fmt.Println(r.URL)
	if r.Method == "CONNECT" {
		handlerHTTPS(w, r)
	} else {
		handlerHTTP(w, r)
	}
}

func CopyHeader(dest, src http.Header) {
	for key, values := range src {
		dest.Set(key, values[0])
	}
}

func main() {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      http.HandlerFunc(handler),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Fatal(server.ListenAndServe())
}
