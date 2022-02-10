package main

import (
	"bytes"
	"fmt"
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
