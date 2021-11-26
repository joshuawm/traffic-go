package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

type requestsPayloadStruct struct {
	ProxyCondition string `json:"proxy_condition"`
}

//Get a json decoder for a  ·11`given requests body
func requestBodyDecoder(r *http.Request) *json.Decoder {
	//Read body to buffer
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("failed to read the request body")
		log.Println("failed to read the request body")
		panic(err)
	}
	//Because golang is a pain in the ass if you read the body then any susequent calls are unable to read the body again....
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body)))
}

func parseRequestBody(r *http.Request) requestsPayloadStruct {
	decoder := requestBodyDecoder(r)

	var requestPayLoad requestsPayloadStruct
	err := decoder.Decode(&requestPayLoad)

	if err != nil {
		panic(err)
	}
	return requestPayLoad
}

func main() {
	http.HandleFunc("/redirect", HandleRequestAndRedirect)
	if err := http.ListenAndServe(":1200", nil); err != nil {
		fmt.Println("failed to start the server")
		panic(err)
	}
}
