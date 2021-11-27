package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//Given a request sned it to the appropriate url
func HandleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	u, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	url := u.Get("url")
	//nothing to write
	// requestPayLoad := parseRequestBody(r)
	// fmt.Println(requestPayLoad.ProxyCondition)
	serveReverseProxy(url, w, r)

}
func serveReverseProxy(target string, w http.ResponseWriter, r *http.Request) {
	//parse the url structure from target
	url, _ := url.Parse(target)

	//create the reserve proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	//Update the header to support the SSL redirection
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host
	r.URL.Path = ""
	//Note that ServeHttp is non blocking and use a go rountine under the hood
	proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/redirect", HandleRequestAndRedirect)
	if err := http.ListenAndServe(":1200", nil); err != nil {
		fmt.Println("failed to start the server")
		panic(err)
	}
}
