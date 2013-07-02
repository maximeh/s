package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
  "runtime"
	"strings"
)

func in_slice(a string, list []string) bool {
	for b := range list {
		if list[b] == a {
			return true
		}
	}
	return false
}

// Shamelessly ripped off of Simon Budig's woof
// http://www.home.unix-ag.org/simon/woof
func find_ip() string {
	ips := [3]string{"192.168.2.0", "198.51.100.0", "203.0.113.0"}
	ip_addr := []string{}
	for ip := range ips {
		serverAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:80", ips[ip]))
		con, _ := net.DialUDP("udp", nil, serverAddr)
		new_ip := con.LocalAddr().String()
		new_host, _, _ := net.SplitHostPort(new_ip)
		// If new_host already in ip_addr, return new_host
		if in_slice(new_host, ip_addr) {
			return new_host
		}
		ip_addr = append(ip_addr, new_host)
	}
	return ip_addr[0]
}

func main() {
	port_i := flag.Int("port", 4242, "port number")
	flag.Parse()
	port := fmt.Sprintf(":%d", *port_i)

	args := flag.Args()

	if len(args) == 0 {
		log.Printf("Nothing to serve.")
		return
	} else if len(args) > 1 {
		log.Printf("We can only serve one element at a time.")
		return
	}

	file_name := filepath.Base(flag.Args()[0])
	file_path, err := filepath.Abs(flag.Args()[0])
	if err != nil {
		log.Printf("Error getting absolute path for %s: %v", file_path, err)
		return
	}

	file_info, err := os.Stat(file_path)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	ip_addr := find_ip()
	log.Printf("Serving %s at http://%s%s", file_path, ip_addr, port)
	url := fmt.Sprintf("http://%s%s", ip_addr, port)
  cmd := exec.Command("xclip", "-i", "-selection", "clipboard")
  if runtime.GOOS == "darwin" {
    cmd = exec.Command("pbcopy")
  }
	cmd.Stdin = strings.NewReader(url)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	handler := http.FileServer(http.Dir(file_path))
	if !file_info.IsDir() {
		handler = nil
		// Easier to use a closure func since we use local variables
		serveFile := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Description", "File Transfer")
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file_name))
			w.Header().Set("Content-Transfer-Encoding", "binary")
			w.Header().Set("Expires", "0")
			w.Header().Set("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
			w.Header().Set("Pragma", "public")
			w.Header().Set("Content-Length", fmt.Sprint("%s", file_info.Size()))
			http.ServeFile(w, r, file_path)
			os.Exit(0)
		})
		http.HandleFunc("/", serveFile)
	}

	err = http.ListenAndServe(port, handler)
	if err != nil {
		log.Printf("Error running web server for static assets: %v", err)
	}
}
