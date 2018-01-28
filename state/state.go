package state

import (
	"net"
  "strings"
  "log"
)

var MyId  string
var MyURL string
var Port  string = ":3000"
var Name  string = "VA"

func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Panicf("Error initializing state: %s", err)
	}

	var url string
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
        url = ipnet.IP.String()
			}
		}
	}
  if url == "" {
    log.Panicf("Error initializing state: couldn't find IP\n")
  }

  strs := strings.Split(url, ".")
  MyId = strs[len(strs) - 1]
	MyURL = url + Port
  log.Printf("My URL is %s", MyURL)
  log.Printf("My id is %s", MyId)
}
