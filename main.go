package main

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type DNSHeader struct {
	id             uint16
	flags          uint16
	numQuestions   uint16
	numAnswers     uint16
	numAuthorities uint16
	numAdditionals uint16
}

type DNSQuestion struct {
	name       []byte
	fieldType  uint16
	fieldClass uint16
}

func main() {

}

func headerToBytes(header DNSHeader) []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, header.id)
	binary.Write(&buf, binary.BigEndian, header.flags)
	binary.Write(&buf, binary.BigEndian, header.numQuestions)
	binary.Write(&buf, binary.BigEndian, header.numAuthorities)
	binary.Write(&buf, binary.BigEndian, header.numAdditionals)

	return buf.Bytes()
}

func encodeDnsName(name string) []byte {
	labels := strings.Split(name, ".")

	var buf bytes.Buffer

	for _, label := range labels {
		if len(label) > 63 {
			label = label[:63]
		}

		buf.WriteByte(byte(len(label)))
		buf.WriteString(label)
	}

	buf.WriteByte(0)
	return buf.Bytes()
}
