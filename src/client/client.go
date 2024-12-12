package main

import (
	"bufio"
	"core/core"
	"fmt"
	"net"
	"net/url"
	"strings"
)

func process(conn net.Conn, depwd core.Password, enpwd core.Password) {
	defer conn.Close()

	proxyServer, err := net.Dial("tcp", "127.0.0.1:8080")

	if err != nil {
		return
	}

	fmt.Println("Connected to proxy server")
	core.EncodeWrite(proxyServer, enpwd, []byte{0x05, 0x01, 0x00})

	_, buf, err := core.DecodeRead(proxyServer, depwd)

	if err != nil || buf[0] != 0x05 || buf[1] != 0x00 {
		return
	}

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading")
		return
	}

	parts := strings.Split(line, " ")

	fmt.Println("parts: ", parts)

	var host string

	if parts[0] == "CONNECT" {
		host, _, err = net.SplitHostPort(parts[1])
		if err != nil {
			fmt.Println("Error splitting host and port:", err)
			return
		}
	} else {
		URL, err := url.Parse(parts[1])

		if err != nil {
			fmt.Println("Parsing Error: ", err)
			return
		}

		host = URL.Hostname()
	}

	hostIP := net.ParseIP(host)

	fmt.Println("host: ", hostIP)

	addrType := "IPv4"
	if hostIP != nil {
		if hostIP.To4() == nil {
			addrType = "IPv6"
		}
	} else {
		addrType = "Domain name"
	}

	buf = []byte{0x05, 0x01, 0x00}
	switch addrType {
	case "IPv4":
		buf = append(buf, 0x01)
		buf = append(buf, hostIP.To4()...)
	case "IPv6":
		buf = append(buf, 0x04)
		buf = append(buf, hostIP.To16()...)
	case "Domain name":
		buf = append(buf, 0x03)
		buf = append(buf, byte(len(host)))
		buf = append(buf, []byte(host)...)
	}

	if parts[0] == "CONNECT" {
		buf = append(buf, byte(443>>8), byte(443&0xff))
	} else {
		buf = append(buf, byte(80>>8), byte(80&0xff))
	}

	fmt.Println("buf: ", buf)

	core.EncodeWrite(proxyServer, enpwd, buf)

	_, buf, err = core.DecodeRead(proxyServer, depwd)
	if err != nil || buf[0] != 0x05 || buf[1] != 0x00 {
		return
	}
	if parts[0] == "CONNECT" {
		conn.Write([](byte)("HTTP/1.1 200 Connection Established\r\n\r\n"))
	}

	defer proxyServer.Close()
	fmt.Println("Start")
	go func() {
		err := core.EncodeCopy(proxyServer, conn, enpwd)
		if err != nil {
			fmt.Println("Error in EncodeCopy: ", err)
			conn.Close()
			proxyServer.Close()
		}
	}()
	core.DecodeCopy(conn, proxyServer, depwd)

	// tmpProxy, err := net.Dial("tcp", "220.181.38.150:80")
	// if err != nil {
	// 	fmt.Println("Error in tmpProxy: ", err)
	// 	return
	// }
	// go io.Copy(tmpProxy, conn)
	// io.Copy(conn, tmpProxy)
	// go io.Copy(proxyServer, conn)
	// io.Copy(conn, proxyServer)
}

func main() {
	decodePassword := core.DecodePassword
	encodePassword := core.EncodePassword
	listen, err := net.Listen("tcp", ":7878")

	if err != nil {
		return
	}

	fmt.Println("Client is running on 7878")

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connect: ", err)
			return
		}
		go process(conn, decodePassword, encodePassword)
	}
}
