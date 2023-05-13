package cmd

import (
	"net"
	"encoding/binary"
	"errors"
	"bufio"
	"os"
	"fmt"
	"log"
	"crypto/sha256"
)

var (
	ErrNotIPv4Address = errors.New("not an IPv4 addres")
)

func IPv4ToInt(ipaddr net.IP) (uint32, error) {
	if ipaddr.To4() == nil {
		return 0, ErrNotIPv4Address
	}
	return binary.BigEndian.Uint32(ipaddr.To4()), nil
}

func IntToIPv4(ipaddr uint32) net.IP {
	ip := make(net.IP, net.IPv4len)

	// Proceed conversion
	binary.BigEndian.PutUint32(ip, ipaddr)

	return ip
}

func ImportDataFile(filePath string) *bufio.Scanner {
	data, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	return bufio.NewScanner(data)
}

func HashChunk(chunk []byte) string {
	h := sha256.New()
	h.Write(chunk)
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func GetChunk(filepath string) []byte {
	file, err := os.Open(filepath)

	if err != nil {
		fmt.Println("cannot able to read the file", err)
		return []byte{}
	}

	defer file.Close()
	// read the first 1024 bytes of the file
	var reader *bufio.Reader = bufio.NewReader(file)

	// Read the file in 4-byte chunks
	var chunkSize int = 1024
	chunk := make([]byte, chunkSize)

	// Read the next chunk
	n, err := reader.Read(chunk)
	if err != nil && err.Error() != "EOF" {
		log.Fatal(err)
	}

	return chunk[:n]
}
