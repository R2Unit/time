package ntp

import (
	"log"
	"net"

	"github.com/r2unit/time/timekeeper"
)

func StartServer(address string, timeKeeper *timekeeper.DriftAwareTimeKeeper) error {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	buffer := make([]byte, 48)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Printf("Failed to read packet: %v", err)
			continue
		}

		if n < 48 {
			log.Printf("Received malformed packet from %s", addr)
			continue
		}

		go handleRequest(conn, addr, buffer, timeKeeper)
	}
}

func handleRequest(conn net.PacketConn, addr net.Addr, req []byte, timeKeeper *timekeeper.DriftAwareTimeKeeper) {
	packet, err := DecodePacket(req)
	if err != nil {
		log.Printf("Failed to decode packet from %s: %v", addr, err)
		return
	}

	resp := &NtpPacket{
		Settings:       0x1C, // LI=0, Version=3, Mode=4 (server)
		Stratum:        1,
		Poll:           4,
		Precision:      -20,
		RootDelay:      0,
		RootDispersion: 0,
		ReferenceID:    0x4C4F434C, // "LOCL"
		RefTimestamp:   timeKeeper.GetCurrentNtpTime(),
		OrigTimestamp:  packet.TxTimestamp,             // Echo the client's transmit timestamp
		RxTimestamp:    timeKeeper.GetCurrentNtpTime(), // Time when the request was received
		TxTimestamp:    timeKeeper.GetCurrentNtpTime(), // Time when the response is sent
	}

	log.Printf("Server Response: TxTimestamp = %d (NTP format)", resp.TxTimestamp)

	respBytes, err := EncodePacket(resp)
	if err != nil {
		log.Printf("Failed to encode response packet: %v", err)
		return
	}

	_, err = conn.WriteTo(respBytes, addr)
	if err != nil {
		log.Printf("Failed to send response to %s: %v", addr, err)
	}
}
