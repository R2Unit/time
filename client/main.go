package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// Seconds from 1900 to 1970
const ntpEpochOffset = 2208988800

func main() {
	server := "localhost:123"

	for {
		ntpTime, err := requestNTPTime(server)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Current NTP time: %v\n", ntpTime)
		}

		time.Sleep(1 * time.Second)
	}
}

func requestNTPTime(server string) (time.Time, error) {
	conn, err := net.Dial("udp", server)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	req := make([]byte, 48)
	req[0] = 0x1B // LI=0, Version=3, Mode=3 (client)
	_, err = conn.Write(req)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to send request: %v", err)
	}

	resp := make([]byte, 48)
	_, err = conn.Read(resp)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to read response: %v", err)
	}

	// Extract transmit timestamp (bytes 40-47)
	secs := binary.BigEndian.Uint32(resp[40:44])  // First 4 bytes: seconds
	fracs := binary.BigEndian.Uint32(resp[44:48]) // Next 4 bytes: fractional seconds

	ntpTime := float64(secs) + float64(fracs)/(1<<32)
	unixTime := ntpTime - ntpEpochOffset

	return time.Unix(int64(unixTime), int64((ntpTime-float64(int64(unixTime)))*1e9)), nil
}
