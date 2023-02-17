package server

import (
	"github.com/golang/glog"
	"net"
	"strconv"
)

func CreateServerUDP(host string, port int, packetChann chan<- []byte) {
	// Resolve the address to listen on
	address := host + ":" + strconv.Itoa(port)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		glog.Error("Error resolving address:", err)
		return
	}
	// Create a UDP listener
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		glog.Error("Error listening:", err)
		return
	}
	glog.Info("Started server UDP receive syslog on port 7000")
	defer conn.Close()
	// Loop to handle incoming packets
	for {
		// Read packet from the UDP connection
		buffer := make([]byte, 2048)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			glog.Error("Error reading packet:", err)
			continue
		}
		// Process the packet
		//fmt.Printf("Received %d bytes from %s: %s\n", n, addr.String(), string(buffer[:n]))
		//Put packet that've received to the channel for processing later
		packetChann <- buffer[:n]
	}
}
