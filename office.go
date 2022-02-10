package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type OfficeEndpoint struct {
	Id                     int      `json:"id"`
	ServiceArea            string   `json:"serviceArea"`
	ServiceAreaDisplayName string   `json:"serviceAreaDisplayName"`
	Urls                   []string `json:"urls"`
	TcpPorts               string   `json:"tcpPorts"`
	Ips                    []string `json:"ips"`
	Notes                  string   `json:"notes"`
}

var oep []OfficeEndpoint

func GetOffice(ref string) []OfficeEndpoint {
	res, err := http.Get(ref)
	if err != nil {
		fmt.Errorf("Could not get the link from microsoft: %v", err)
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&oep)
	return oep
}

func GetOfficeUrls(w http.ResponseWriter, r *http.Request) {
	var ips []string
	eps := GetOffice(officeRef)
	w.Header().Add("Content-Type", "Application/json")
	for _, v := range eps {
		if v.Urls != nil {
			for _, str := range v.Urls {
				ips = append(ips, str)
			}
		}
	}
	fmt.Fprint(w, ips)
}

// func GetOffice(w http.ResponseWriter, r *http.Request) {
// 	eps := GetOffice(officeRef)
// 	w.Header().Add("Content-Type", "Application/json")
// 	json.NewEncoder(w).Encode(eps)
// }

func GetOfficeIpv4(w http.ResponseWriter, r *http.Request) {
	var ips []string
	eps := GetOffice(officeRef)
	fmt.Printf("this is the eps : %v", eps)
	w.Header().Add("Content-Type", "Application/json")
	for _, v := range eps {
		if v.Ips != nil {
			for _, str := range v.Ips {
				if !strings.Contains(str, ":") {
					ips = append(ips, str)
				}

			}
		}
	}
	fmt.Println(len(ips))
	fmt.Fprint(w, ips)
}

func GetOfficeCisco(w http.ResponseWriter, r *http.Request) {
	var ips []string
	var list = make(map[string][]string)
	eps := GetOffice(officeRef)

	w.Header().Add("Content-Type", "Application/json")
	for _, v := range eps {

		name := fmt.Sprintf("O365.%v.%v", v.Id, v.ServiceArea)
		if v.Ips != nil {
			for _, str := range v.Ips {
				if !strings.Contains(str, ":") { // filter ipv6
					_, ipNet, _ := net.ParseCIDR(str)
					mask := ipv4MaskString(ipNet.Mask)
					ips = append(ips, fmt.Sprintf("%v/%v", ipNet.IP.String(), mask))
				}
			}
		}
		if len(ips) > 0 {
			list[name] = ips
		}
		fmt.Println((list))
		ips = ips[:0]

	}
	str := parseASArule(list)

	w.Write([]byte(str.String()))
}
