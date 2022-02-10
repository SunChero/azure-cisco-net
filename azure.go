package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

var azep AzureEndpoint

type AzureEndpoint struct {
	Values []struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		Properties struct {
			Region          string   `json:"region"`
			RegionId        int      `json:"regionId"`
			AddressPrefixes []string `json:"addressPrefixes"`
		} `json:"properties"`
	} `json:"values"`
}

func GetAzure(ref string) AzureEndpoint {
	res, err := http.Get(ref)
	if err != nil {
		fmt.Errorf("Could not get the link from microsoft: %v", err)
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&azep)
	return azep
}

// http.HandleFunc("/azure", func(w http.ResponseWriter, r *http.Request) {
// 	aps := GetAzure(azureRef)
// 	w.Header().Add("Content-Type", "Application/json")
// 	json.NewEncoder(w).Encode(aps)
// })
func GetAzureIpv4(w http.ResponseWriter, r *http.Request) {
	var ips []string
	eps := GetAzure(azureRef)
	w.Header().Add("Content-Type", "Application/json")
	for _, v := range eps.Values {
		if v.Properties.AddressPrefixes != nil {
			for _, str := range v.Properties.AddressPrefixes {
				if !strings.Contains(str, ":") { // filter ipv6
					_, ipNet, _ := net.ParseCIDR(str)
					mask := ipv4MaskString(ipNet.Mask)
					ips = append(ips, fmt.Sprintf("%v - %v", ipNet.IP.String(), mask))
				}
			}
		}
	}
	fmt.Fprintf(w, "%v \n", ips)
}

func GetAzureCisco(w http.ResponseWriter, r *http.Request) {

	var ips []string
	var list = make(map[string][]string)
	eps := GetAzure(azureRef)
	w.Header().Add("Content-Type", "Application/json")

	for _, v := range eps.Values {
		//access to rules name
		name := fmt.Sprintf("AZ.%v.%v", v.Properties.RegionId, v.Name)
		if v.Properties.AddressPrefixes != nil {
			for _, str := range v.Properties.AddressPrefixes {
				if !strings.Contains(str, ":") { // filter ipv6
					_, ipNet, _ := net.ParseCIDR(str)
					mask := ipv4MaskString(ipNet.Mask)
					ips = append(ips, fmt.Sprintf("%v/%v", ipNet.IP.String(), mask))
				}
			}
		}
		list[name] = ips
		ips = ips[:0]

	}
	key, ok := r.URL.Query()["search"]
	if ok {
		for k, _ := range list {
			if !strings.Contains(strings.ToUpper(k), strings.ToUpper(key[0])) {
				delete(list, k)
			}
		}
	}
	str := parseASArule(list)

	w.Write([]byte(str.String()))
}
