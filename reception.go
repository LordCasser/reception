package reception

import (
	"github.com/LordCasser/reception/utils"
	"inet.af/tcpproxy"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type Reception struct {
	port     int
	sslPort  int
	mux      tcpproxy.Proxy
	Redirect []string
	Switch   map[string]*httputil.ReverseProxy
}

func New() *Reception {
	return &Reception{Switch: make(map[string]*httputil.ReverseProxy), mux: tcpproxy.Proxy{}, port: 80, sslPort: 443}
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
	to, err := url.Parse(strings.TrimSpace(transfer))
	if err != nil {
		log.Panicln("wrong host url!\n", err)
	}
	in, err := url.Parse(strings.TrimSpace(host))
	if err != nil {
		log.Panicln("wrong host url!\n", err)
	}
	if to.Scheme == "https" {
		x.mux.AddSNIRoute(":"+strconv.Itoa(x.sslPort), in.Hostname(), tcpproxy.To(to.Host))
		x.Redirect = append(x.Redirect, in.Hostname())
	} else {
		x.Switch[host] = httputil.NewSingleHostReverseProxy(to)
	}
	return nil
}

func (x *Reception) Serve() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Host)
		if proxy, ok := x.Switch[r.Host]; ok {
			proxy.ServeHTTP(w, r)
		} else if utils.ContainsInSlice(x.Redirect, r.Host) {
			utils.Redirect(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	go func() {
		err := x.mux.Run()
		if err != nil {
			log.Println(err)
		}
	}()
	err := http.ListenAndServe(":"+strconv.Itoa(x.port), nil)
	if err != nil {
		log.Println(err)
	}
}
