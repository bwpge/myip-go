package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	flag "github.com/spf13/pflag"
)

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
	port := flag.UintP("port", "p", 8080, "the port to listen on")
	iface := flag.StringP("interface", "i", "", "optional network interface address to use")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *iface, *port)

	http.HandleFunc("/", handler)
	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
