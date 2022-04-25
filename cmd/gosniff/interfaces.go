package gosniff

import (
	"fmt"

	"github.com/google/gopacket/pcap"
)

// GetInterfaces returns all host interfaces in string format
func GetInterfaces() (interfaces []pcap.Interface, err error) {
	ifaces, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Println("Error: No host interfaces")
		return interfaces, err
	}

	for _, i := range ifaces {
		// if len(i.Addresses) > 0 {
		interfaces = append(interfaces, i)
		// }
	}

	return interfaces, nil
}
