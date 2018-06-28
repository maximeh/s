package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var HEADERS = map[string]string{
	"Content-Description":       "File Transfer",
	"Content-Type":              "application/octet-stream",
	"Content-Transfer-Encoding": "binary",
	"Cache-Control":             "private",
	"Pragma":                    "private",
	"Expires":                   "Mon, 26 Jul 1997 05:00:00 GMT",
}

type DownloadFile struct {
	name  string
	path  string
	size  string
	count int
}

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

func serve(dl_file DownloadFile) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		for key, value := range HEADERS {
			w.Header().Set(key, value)
		}

		cd_value := fmt.Sprintf("attachment; filename=%s", dl_file.name)
		w.Header().Set("Content-Disposition", cd_value)
		w.Header().Set("Content-Length", dl_file.size)

		http.ServeFile(w, r, dl_file.path)
		dl_file.count--
		if dl_file.count == 0 {
			os.Exit(0)
		}
		log.Printf("Download left: %d\n", dl_file.count)
	}
}

func main() {
	usage := `
Usage:
  s [options] <path>
  s -h | --help
  s -v | --version

Options:
  -h --help            Show this screen.
  -v --version         Show version.
  -c --count=<count>   Port to use [default: 1].
  -p --port=<port>     Port to use [default: 4242].
`
	arguments, _ := docopt.Parse(usage, nil, true, "s 1.0", false)
	var dl_file DownloadFile
	var err error

	port := fmt.Sprintf(":%s", arguments["--port"])
	dl_file.count, err = strconv.Atoi(arguments["--count"].(string))
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	dl_file.name = filepath.Base(arguments["<path>"].(string))
	dl_file.path, err = filepath.Abs(arguments["<path>"].(string))
	if err != nil {
		log.Printf("Error getting absolute path for %s: %v", dl_file.path, err)
		return
	}

	file_info, err := os.Stat(dl_file.path)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	dl_file.size = strconv.FormatInt(file_info.Size(), 10)

	ip_addr := find_ip()
	url := fmt.Sprintf("http://%s%s/%s", ip_addr, port, dl_file.name)
	log.Printf("Serving %s at %s", dl_file.path, url)
	cmd := exec.Command("xclip", "-i", "-selection", "clipboard")
	cmd.Env = append(cmd.Env, "DISPLAY=:0")
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("pbcopy")
	}
	cmd.Stdin = strings.NewReader(url)
	err = cmd.Run()
	if err != nil {
		log.Print("Note: The URL could not be copied in your clipboard.")
	}

	handler := http.FileServer(http.Dir(dl_file.path))
	if !file_info.IsDir() {
		handler = nil
		serveFile := http.HandlerFunc(serve(dl_file))
		http.HandleFunc("/", serveFile)
	}

	err = http.ListenAndServe(port, handler)
	if err != nil {
		log.Printf("Error running web server for static assets: %v", err)
	}
}
