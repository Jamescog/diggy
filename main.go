package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strings"
)

const (
	Class int = 1
	Type  int = 1
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

type DNSRecord struct {
	name       []byte
	fieldType  uint16
	fieldClass uint16
	feildTTL   uint32
	rdLength   uint16
	rData      []byte
}

func main() {
	domain := "www.example.com"

	query, err := buildQuery(domain, Type)

	if err != nil {
		fmt.Print("Error occured")
	}
	res, err := sendQuery(query, "8.8.8.8:53")
	header, err := parseHeader(res)

	if err != nil {
		fmt.Printf("Error occured: %v", err)
		return
	}

	fmt.Printf("Header: %v", header)

}

func headerToBytes(header DNSHeader) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, header.id)
	binary.Write(&buf, binary.BigEndian, header.flags)
	binary.Write(&buf, binary.BigEndian, header.numQuestions)
	binary.Write(&buf, binary.BigEndian, header.numAnswers)
	binary.Write(&buf, binary.BigEndian, header.numAuthorities)
	binary.Write(&buf, binary.BigEndian, header.numAdditionals)
	return buf.Bytes()
}

func questionToBytes(question DNSQuestion) ([]byte, error) {
	var buf bytes.Buffer

	_, err := buf.Write(question.name)
	if err != nil {
		return nil, fmt.Errorf("failed to write question name: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, question.fieldType)

	if err != nil {
		return nil, fmt.Errorf("failed to write question type: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, question.fieldClass)

	if err != nil {
		return nil, fmt.Errorf("failed to write question class: %w", err)
	}

	return buf.Bytes(), nil
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

func genHeader() DNSHeader {
	id := rand.Intn(65535)
	recurrsionDesired := 1 << 8
	header := DNSHeader{}
	header.id = uint16(id)
	header.flags = uint16(recurrsionDesired)
	header.numQuestions = uint16(1)
	header.numAnswers = 0
	header.numAdditionals = 0
	return header
}
func genQuestion(name string, recordType int) DNSQuestion {
	nameByte := encodeDnsName(name)
	question := DNSQuestion{}
	question.name = nameByte
	question.fieldType = uint16(recordType)
	question.fieldClass = uint16(Class)
	return question
}

func buildQuery(domainName string, recordType int) ([]byte, error) {

	var buf bytes.Buffer
	header := genHeader()
	question := genQuestion(domainName, recordType)
	headerByte := headerToBytes(header)
	questionByte, err := questionToBytes(question)

	if err != nil {
		return nil, fmt.Errorf("failed to generating question bytes: %w", err)
	}

	_, err = buf.Write(headerByte)

	if err != nil {
		return nil, fmt.Errorf("failed to write header bytes: %w", err)
	}

	_, err = buf.Write(questionByte)

	if err != nil {
		return nil, fmt.Errorf("failed to write question bytes: %w", err)
	}

	return buf.Bytes(), nil

}

func parseHeader(data []byte) (DNSHeader, error) {
	if len(data) < 12 {
		return DNSHeader{}, fmt.Errorf("response is too short: %d", len(data))
	}

	reader := bytes.NewReader(data)
	var header DNSHeader

	err := binary.Read(reader, binary.BigEndian, &header.id)

	if err != nil {
		return DNSHeader{}, fmt.Errorf("Error reading header: %v", err)
	}

	err = binary.Read(reader, binary.BigEndian, &header.flags)

	if err != nil {
		return DNSHeader{}, fmt.Errorf("Error reading flags: %v", err)
	}

	err = binary.Read(reader, binary.BigEndian, &header.numQuestions)

	if err != nil {
		return DNSHeader{}, fmt.Errorf("Error reading questions: %v", err)
	}

	err = binary.Read(reader, binary.BigEndian, &header.numAnswers)

	if err != nil {
		return DNSHeader{}, fmt.Errorf("Error reading answers: %v", err)
	}

	err = binary.Read(reader, binary.BigEndian, &header.numAuthorities)

	if err != nil {
		return DNSHeader{}, fmt.Errorf("Error reading authorities: %v", err)
	}

	err = binary.Read(reader, binary.BigEndian, &header.numAdditionals)

	if err != nil {
		return DNSHeader{}, fmt.Errorf("Error reading addtionals: %v", err)
	}

	return header, nil

}

func sendQuery(query []byte, server string) ([]byte, error) {

	addr, err := net.ResolveUDPAddr("udp", server)

	if err != nil {
		return nil, fmt.Errorf("failded to resolve server address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		return nil, fmt.Errorf("failded to dial UPD: %w", err)
	}

	defer conn.Close()

	_, err = conn.Write(query)

	if err != nil {
		return nil, fmt.Errorf("failed to send query: %w", err)
	}

	response := make([]byte, 512)
	n, _, err := conn.ReadFromUDP(response)

	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return response[:n], nil
}
