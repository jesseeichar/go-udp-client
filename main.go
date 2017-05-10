package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func checkError(desc string, err error) {
	if err != nil {
		fmt.Println(desc+"Error: ", err)
		os.Exit(0)
	}
}

func main() {
	host := flag.String("host", "127.0.0.1", "udp host of server to contact to ")
	port := flag.Int("port", 90, "udp port of server to contact to")
	data := flag.String("data", "Test Data from Client", "data to send to server")

	address := fmt.Sprintf("%s:%d", *host, *port)

	fmt.Printf("Connecting to %s\n", address)

	conn, err := net.Dial("udp", address)
	checkError("open connection", err)
	defer conn.Close()

	go func() {
		fmt.Printf("Writing test data, timeout 10 seconds\n")
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_, err = conn.Write([]byte(*data))
		checkError("Writing data", err)
	}()

	serverAddr, err := net.ResolveUDPAddr("udp", conn.LocalAddr().String())
	checkError("Resolve udp address", err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	checkError("Start udp listener", err)
	defer serverConn.Close()
	serverConn.SetReadDeadline(time.Now().Add(10 * time.Second))
	fmt.Printf("Reading response data, timeout 10 seconds\n")
	buf := make([]byte, 1024)
	n, _, err := serverConn.ReadFromUDP(buf)
	checkError("Error reading data from udp", err)

	fmt.Printf("Data read: %s\n", string(buf[0:n]))
}
