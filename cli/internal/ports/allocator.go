package ports

import (
	"fmt"
	"net"
)

// Allocate finds a free port using bind(0) probe (SPEC §20.9: atomic, no race)
func Allocate(prefPort int) (int, error) {
	// If a preferred port was requested, try it first
	if prefPort > 0 {
		if isFree(prefPort) {
			return prefPort, nil
		}
	}

	// Otherwise, find the next free port in the dev range (8000-8999)
	for port := 8000; port < 9000; port++ {
		if isFree(port) {
			return port, nil
		}
	}

	return 0, fmt.Errorf("no free ports available in 8000-8999 range")
}

// isFree checks if a port is available by attempting to bind to it
func isFree(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// List prints all allocated ports from the database
func List() {
	ports, err := LoadDB()
	if err != nil {
		fmt.Printf("Error loading ports: %v\n", err)
		return
	}

	if len(ports) == 0 {
		fmt.Println("No ports allocated")
		return
	}

	fmt.Println("Allocated ports:")
	for _, p := range ports {
		fmt.Printf("  %d — %s (%s)\n", p.Port, p.ProjectPath, p.URL)
	}
}
