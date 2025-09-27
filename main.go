package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type IPResponse struct {
	RemoteAddr   string `json:"ip"`
	ForwardedFor string `json:"forwardedFor,omitempty"`
	RealIP       string `json:"realIp,omitempty"`
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
	resp := IPResponse{}

	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		resp.RealIP = normalize(ip)
	}

	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		resp.ForwardedFor = normalize(strings.Split(ip, ",")[0])
	}

	resp.RemoteAddr = normalize(r.RemoteAddr)
	return resp
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
	fmt.Fprintln(w, "IP:", result.RemoteAddr)
	if result.ForwardedFor != "" {
		fmt.Fprintln(w, "Forwarded for:", result.ForwardedFor)
	}
	if result.RealIP != "" {
		fmt.Fprintln(w, "Real IP:", result.RealIP)
	}
}

func main() {
	port := uint64(8080)
	if len(os.Args) == 2 {
		p, err := strconv.ParseUint(os.Args[1], 10, 16)
		if err != nil || p == 0 {
			log.Fatal("invalid argument")
		}
		port = p
	} else if len(os.Args) != 1 {
		log.Fatal("invalid usage")
	}
	addr := fmt.Sprintf(":%d", port)

	http.HandleFunc("/", handler)
	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
