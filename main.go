package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	log "github.com/bwpge/systemdlog"
)

var confPath = "/etc/myip/myip.json"

type config struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}

type IPResponse struct {
	IP string `json:"ip"`
}

func normalize(s string) string {
	host, _, err := net.SplitHostPort(s)
	if err != nil {
		return s
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return s
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		ipv6 := ip.To16().String()
		// usually get this when talking to localhost vs 127.0.0.1
		if ipv6 == "::1" {
			return "127.0.0.1"
		} else {
			return ipv6
		}
	} else {
		return ipv4.String()
	}
}

func extractIP(r *http.Request) IPResponse {
	result := IPResponse{}

	// a proxy or load balancer is usually going to set this
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		result.IP = normalize(strings.Split(ip, ",")[0])
		return result
	}

	// this is non-standard, but is still used enough
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		result.IP = normalize(ip)
		return result
	}

	// fallback to remote address
	result.IP = normalize(r.RemoteAddr)
	return result
}

func handler(w http.ResponseWriter, r *http.Request) {
	result := extractIP(r)
	accept := r.Header.Get("Accept")

	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, result.IP)
}

func main() {
	conf := config{}
	confData, err := os.ReadFile(confPath)
	if err == nil {
		if err := json.Unmarshal(confData, &conf); err != nil {
			log.Fatalf("failed to load config: %s", err)
		}
	}

	port := flag.UintP("port", "p", 0, "the port to listen on")
	host := flag.StringP("host", "h", "", "IP address or hostname to bind")
	flag.Parse()

	if *host != "" {
		conf.Host = *host
	}
	if *port != 0 {
		conf.Port = *port
	}
	if conf.Port == 0 {
		log.Fatal("port must be provided through config or command line arguments")
	}

	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	http.HandleFunc("/", handler)
	log.Infof("server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
