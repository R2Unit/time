package ntp

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type NtpPacket struct {
	Settings       uint8
	Stratum        uint8
	Poll           int8
	Precision      int8
	RootDelay      uint32
	RootDispersion uint32
	ReferenceID    uint32
	RefTimestamp   uint64
	OrigTimestamp  uint64
	RxTimestamp    uint64
	TxTimestamp    uint64
}

func DecodePacket(data []byte) (*NtpPacket, error) {
	if len(data) < 48 {
		return nil, errors.New("invalid NTP packet size")
	}
	packet := &NtpPacket{}
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, packet)
	return packet, err
}

func EncodePacket(packet *NtpPacket) ([]byte, error) {
	buf := make([]byte, 48)
	err := binary.Write(bytes.NewBuffer(buf), binary.BigEndian, packet)
	return buf, err
}
