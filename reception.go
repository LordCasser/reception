package reception

import (
	"inet.af/tcpproxy"
	"log"
	"strconv"
)

type Reception struct {
	port    int
	sslPort int
	mux     tcpproxy.Proxy
}

func New() *Reception {
	return &Reception{mux: tcpproxy.Proxy{}, port: 80, sslPort: 443}
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
func (x *Reception) AddSwitch(host string, transfer string, ssl bool) error {
	log.Println("host:", host, "trans:", transfer)
	if ssl {
		x.mux.AddSNIRoute(":"+strconv.Itoa(x.sslPort), host, tcpproxy.To(transfer))
	} else {
		x.mux.AddHTTPHostRoute(":"+strconv.Itoa(x.port), host, tcpproxy.To(transfer))
	}

	return nil
}

func (x *Reception) Serve() {
	err := x.mux.Run()

	if err != nil {
		log.Println(err)
	}
}
