package tuya

import (
	"fmt"
	"log"
	"net"
)

func UDPListener() {
	cnx, err := net.ListenPacket("udp", ":6667")
	if err != nil {
		log.Fatal("UDP Listener failed:", err)
	}
	for {
		buffer := make([]byte, 1024)
		n, _, err := cnx.ReadFrom(buffer)
		buffer = buffer[:n]
		if err != nil || len(buffer) < 16 {
			continue
		}

		m := NewMessageUdp()
		msg, err := m.Parse(buffer)
		if err == nil {
			fmt.Printf("%#v\n", msg)
		}
	}
}
