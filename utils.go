package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//utils ...
func ipv4MaskString(m []byte) string {
	if len(m) != 4 {
		panic("ipv4Mask: len must be 4 bytes")
	}

	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
}

func parseASArule(list map[string][]string) bytes.Buffer {
	var out bytes.Buffer
	for key, value := range list {
		out.WriteString(fmt.Sprintf("no object-group network %v \n", key))
		out.WriteString(fmt.Sprintf("object-group network %v \n", key))
		for _, v := range value {
			parts := strings.Split(v, "/")
			if parts[1] == "255.255.255.255" {
				out.WriteString(fmt.Sprintf("network-object host %v  %v\n", parts[0], parts[1]))
			} else {
				out.WriteString(fmt.Sprintf("network-object %v  %v\n", parts[0], parts[1]))
			}
		}

	}
	return out

}

func parsePAC(list []string) string {
	var out bytes.Buffer

	for key, val := range list {
		if strings.Contains(val, "*") {
			if key == (len(list) - 1) {
				out.WriteString(fmt.Sprintf("\tshExpMatch(url, '%v') \n", val))
			} else {
				out.WriteString(fmt.Sprintf("\tshExpMatch(url, '%v') || \n", val))
			}
		} else {
			if key == (len(list) - 1) {
				out.WriteString(fmt.Sprintf("\tdnsDomainIs(host,%v) \n", val))
			} else {
				out.WriteString(fmt.Sprintf("\tdnsDomainIs(host,%v) || \n", val))
			}
		}
	}
	f, _ := filepath.Abs("./proxy.pac")
	pac, err := os.ReadFile(f)
	if err != nil {
		fmt.Printf("file doesnt exist: %v", err)
	}
	leftsplit := strings.Split(string(pac), "//StartSync")
	str1 := leftsplit[0]
	rightsplit := strings.Split(leftsplit[1], "//EndSync")
	str2 := rightsplit[1]
	ret := str1 + fmt.Sprintln("\t//StartSync") + out.String() + fmt.Sprintln("\t//EndSync") + str2
	//fmt.Println(ret)
	return ret

}

func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}
