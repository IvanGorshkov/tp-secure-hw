package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

func handlerHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Proxy-Connection")
	proxyResponse, err := RequestViaProxy(r)
	if err != nil {
		fmt.Println(err)
	}
	defer proxyResponse.Body.Close()
	CopyResponse(proxyResponse, w)
}

func RequestViaProxy(r *http.Request) (*http.Response, error) {
	proxyResponse, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	return proxyResponse, nil
}

func CopyResponse(proxyResponse *http.Response, w http.ResponseWriter) {
	for header, values := range proxyResponse.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.WriteHeader(proxyResponse.StatusCode)
	_, err := io.Copy(w, proxyResponse.Body)
	if err != nil {
		fmt.Println(err)
	}
}

func handlerHTTPS(w http.ResponseWriter, r *http.Request) {

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
