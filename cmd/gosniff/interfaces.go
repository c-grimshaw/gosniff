package gosniff

import (
	"fmt"
	"net"
)

// GetInterfaces returns all host interfaces in string format
func GetInterfaces() ([]string, error) {
	interfaces := []string{}

	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error: No host interfaces")
		return interfaces, err
	}

	for _, i := range ifaces {
		interfaces = append(interfaces, i.Name)
	}

	return interfaces, nil
}
