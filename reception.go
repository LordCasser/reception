package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type Reception struct {
	port int
	//SSLPort int
	Switch map[string]*httputil.ReverseProxy
}

func New() *Reception {
	//if listenPort <= 0 || listenPort >= 65535 {
	//	log.Panicln("wrong listen port")
	//}
	return &Reception{Switch: make(map[string]*httputil.ReverseProxy)}
}
func (x *Reception) SetPort(port int) {
	//func (x *Reception) SetPort(port int, sslPort int) {
	//if (port <= 0 || port >= 65535) || (sslPort <= 0 || sslPort >= 65535) || port == sslPort {
	if port <= 0 || port >= 65535 {
		log.Panicln("wrong listen port")
	}
	x.port = port
	//x.SSLPort = sslPort
}

// AddSwitch @url and transfer string will be parsed here, url example: https://www.example.com
func (x *Reception) AddSwitch(host string, transfer string) error {
	u, err := url.Parse(strings.TrimSpace(transfer))
	if err != nil {
		log.Panicln("wrong host url!\n", err)
	}
	x.Switch[host] = httputil.NewSingleHostReverseProxy(u)
	return nil
}

func (x *Reception) Serve() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if proxy, ok := x.Switch[r.Host]; ok {
			proxy.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	if x.port == 0 {
		x.port = 80
	}
	err := http.ListenAndServe(":"+strconv.Itoa(x.port), nil)
	if err != nil {
		log.Println(err)
	}
}

//func main() {
//	rec := New()
//	_ = rec.AddSwitch("localhost:8080", "http://127.0.0.1:8081")
//	_ = rec.AddSwitch("127.0.0.1:8080", "http://127.0.0.1:8082")
//	rec.Serve()
//}
