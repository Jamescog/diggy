package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
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

func main() {
	domain := "www.example.com"

	query, err := buildQuery(domain, Type)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("Query: %v", query)
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

func buildQuery(domainName string, recordType int) ([]byte, error) {

	var buf bytes.Buffer

	nameByte := encodeDnsName(domainName)
	id := rand.Intn(65535)
	recurrsionDesired := 1 << 8
	header := DNSHeader{}
	header.id = uint16(id)
	header.flags = uint16(recurrsionDesired)
	header.numQuestions = uint16(1)
	question := DNSQuestion{}
	question.name = nameByte
	question.fieldType = uint16(recordType)
	question.fieldClass = uint16(Class)

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
