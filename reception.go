package reception

import (
	"github.com/kevinpollet/tlsmux"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type Reception struct {
	port    int
	sslPort int
	mux     tlsmux.Mux
	Switch  map[string]*httputil.ReverseProxy
}

func New() *Reception {
	return &Reception{Switch: make(map[string]*httputil.ReverseProxy), mux: tlsmux.Mux{}, port: 80, sslPort: 443}
}

func (x *Reception) SetPort(port int) {
	if (port <= 0 || port >= 65535) || port == x.sslPort {
		log.Panicln("wrong listen port")
	}
	x.port = port

}
func (x *Reception) SetSPort(sport int) {
	if (sport <= 0 || sport >= 65535) || sport == x.port {
		log.Panicln("wrong listen ssl port")
	}
	x.sslPort = sport

}

// AddSwitch @url and transfer string will be parsed here, host example: example.com ,transfer example: https://www.example.com
func (x *Reception) AddSwitch(host string, transfer string) error {
	u, err := url.Parse(strings.TrimSpace(transfer))
	if err != nil {
		log.Panicln("wrong host url!\n", err)
	}
	if u.Scheme == "https" {
		x.mux.Handle(host, tlsmux.ProxyHandler{Addr: u.String()})
	} else {
		x.Switch[host] = httputil.NewSingleHostReverseProxy(u)
	}
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
	l, err := net.Listen("tcp", ":"+strconv.Itoa(x.sslPort))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		err := x.mux.Serve(l)
		if err != nil {
			log.Println(err)
		}
	}()
	err = http.ListenAndServe(":"+strconv.Itoa(x.port), nil)
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
