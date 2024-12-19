package main

import (
	"core/internal/core"
	"fmt"
	"net"
)

func process(conn net.Conn, depwd core.Password, enpwd core.Password) {
	defer conn.Close()

	fmt.Println("New connection")

	n, buf, err := core.DecodeRead(conn, depwd)

	//fmt.Println(buf)

	if err != nil || n < 3 || buf[0] != 0x05 {
		return
	}

	core.EncodeWrite(conn, enpwd, []byte{0x05, 0x00})

	n, buf, err = core.DecodeRead(conn, depwd)

	if err != nil || n < 4 || buf[0] != 0x05 || buf[1] != 0x01 {
		return
	}

	fmt.Println("get addr message")

	//SOCKS5 解目标地址
	var dIP []byte
	switch buf[3] {
	case 0x01:
		dIP = buf[4 : 4+net.IPv4len]
	case 0x03:
		ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:n-2]))
		if err != nil {
			fmt.Println("Error resolving domain name: ", err)
			return
		}
		dIP = ipAddr.IP
	case 0x04:
		dIP = buf[4 : 4+net.IPv6len]
	default:
		return
	}

	fmt.Println("solve dst addr")

	dPort := buf[n-2 : n]
	dstAddr := &net.TCPAddr{
		IP:   dIP,
		Port: int(dPort[0])<<8 + int(dPort[1]),
	}

	dstServer, err := net.Dial("tcp", dstAddr.String())
	if err != nil {
		fmt.Println("Error connecting to destination server: ", err)
		return
	} else {
		defer dstServer.Close()

		core.EncodeWrite(conn, enpwd, []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	}

	fmt.Println("dstAddr: ", dstAddr)

	fmt.Println("Connected to dst server")
	//建立连接成功

	// go io.Copy(dstServer, conn)
	// io.Copy(conn, dstServer)

	go func() {
		err = core.DecodeCopy(dstServer, conn, depwd)
		if err != nil {
			fmt.Println("Error in DecodeCopy: ", err)
			conn.Close()
			dstServer.Close()
		}
	}()
	core.EncodeCopy(conn, dstServer, enpwd)
}
func main() {
	decodePassword := core.DecodePassword
	encodePassword := core.EncodePassword
	listen, err := net.Listen("tcp", ":7879")

	if err != nil {
		return
	}

	fmt.Println("Server is running on 7879")

	for {
		conn, _ := listen.Accept()
		go process(conn, decodePassword, encodePassword)
	}

}
