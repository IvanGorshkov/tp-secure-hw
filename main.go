package main

import (
	"crypto/tls"
	"fmt"
	"io"
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

	for key, values := range w.Header() {
		resp.Header.Set(key, values[0])
	}

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
	_, err := io.Copy(writeCloser, readCloser)
	if err != nil {
		fmt.Println(err)
	}
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

func main() {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      http.HandlerFunc(handler),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	fmt.Println(server.ListenAndServe())
}
