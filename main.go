package main

import (
	"log"
	"net/http"
)

var officeRef = "https://endpoints.office.com/endpoints/worldwide?clientrequestid=b10c5ed1-bad1-445f-b386-b919946339a7"
var azureRef = "https://download.microsoft.com/download/7/1/D/71D86715-5596-4529-9B13-DA13A5DE5B63/ServiceTags_Public_20220207.json"

func main() {

	http.HandleFunc("/office/ipv4", GetOfficeIpv4)
	http.HandleFunc("/office/ipv4/cisco", GetOfficeCisco)
	http.HandleFunc("/office/urls", GetOfficeUrls)
	http.HandleFunc("/azure/ipv4", GetAzureIpv4)
	http.HandleFunc("/azure/ipv4/cisco", GetAzureCisco)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
