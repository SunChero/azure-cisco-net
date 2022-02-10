package main

import (
	"log"
	"net/http"
)

var officeRef = "https://endpoints.office.com/endpoints/worldwide?clientrequestid=b10c5ed1-bad1-445f-b386-b919946339a7"
var azureRef = "https://download.microsoft.com/download/7/1/D/71D86715-5596-4529-9B13-DA13A5DE5B63/ServiceTags_Public_20220207.json"

func main() {
	fs := http.FileServer(http.Dir("./frontend/build"))
	http.Handle("/", fs)
	http.HandleFunc("/o/ips", GetOfficeIpv4)
	http.HandleFunc("/o/asa", GetOfficeCisco)
	http.HandleFunc("/o/pac", GetOfficeUrls)
	http.HandleFunc("/z/ips", GetAzureIpv4)
	http.HandleFunc("/z/asa", GetAzureCisco)
	log.Println("Server starting on :443")
	log.Println("Current Endpoints")
	log.Println("/o/ips")
	log.Println("/o/ips/asa")
	log.Println("/z/ip")
	log.Println("/z/asa")

	if err := http.ListenAndServe(":443", nil); err != nil {
		log.Fatalf("Cant startup Server: %v", err)
	}

}
