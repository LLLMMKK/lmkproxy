package main

import (
	"core/core"
	"fmt"
	"io"
	"net"
)

func decodeRead(conn net.Conn, depwd core.Password) (int, []byte, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	core.Decode(depwd, buf)
	return n, buf, err
}

func encodeRead(conn net.Conn, enpwd core.Password) (int, []byte, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	core.Encode(enpwd, buf)
	return n, buf, err
}

func encodeWrite(conn net.Conn, enpwd core.Password, buf []byte) (int, error) {
	core.Encode(enpwd, buf)
	n, err := conn.Write(buf)
	return n, err
}

// func decodeWrite(conn net.Conn, depwd core.Password, buf []byte) (int, error) {
// 	core.Decode(depwd, buf)
// 	n, err := conn.Write(buf)
// 	return n, err
// }

func decodeCopy(dst net.Conn, src net.Conn, depwd core.Password) error {
	for {
		n, buf, err := decodeRead(src, depwd)
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		if n > 0 {
			_, err = dst.Write(buf[:n])
			if err != nil {
				return err
			}
		}
	}
}

func encodeCopy(dst net.Conn, src net.Conn, enpwd core.Password) error {

	for {

		n, buf, err := encodeRead(src, enpwd)
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		if n > 0 {
			_, err = dst.Write(buf[:n])
			if err != nil {
				return err
			}
		}
	}
}

func process(conn net.Conn, depwd core.Password, enpwd core.Password) {
	defer conn.Close()

	fmt.Println("New connection")

	n, buf, err := decodeRead(conn, depwd)

	//fmt.Println(buf)

	if err != nil || n < 3 || buf[0] != 0x05 {
		return
	}

	fmt.Println("Step1")

	encodeWrite(conn, enpwd, []byte{0x05, 0x00})

	n, buf, err = decodeRead(conn, depwd)

	if err != nil || n < 4 || buf[0] != 0x05 || buf[1] != 0x01 {
		return
	}

	fmt.Println("Step2")

	//SOCKS5 解目标地址
	var dIP []byte
	switch buf[3] {
	case 0x01:
		dIP = buf[4 : 4+net.IPv4len]
	case 0x03:
		ipAddr, _ := net.ResolveIPAddr("ip", string(buf[5:n-2]))
		dIP = ipAddr.IP
	case 0x04:
		dIP = buf[4 : 4+net.IPv6len]
	default:
		return
	}

	fmt.Println("Step2Over")

	dPort := buf[n-2 : n]
	dstAddr := &net.TCPAddr{
		IP:   dIP,
		Port: int(dPort[0])<<8 + int(dPort[1]),
	}

	dstServer, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		return
	} else {
		defer dstServer.Close()
		dstServer.SetLinger(0)

		encodeWrite(conn, enpwd, []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	}

	fmt.Println(dstAddr)

	fmt.Println("Connected to dst server")
	//建立连接成功

	go func() {
		err = decodeCopy(dstServer, conn, depwd)
		if err != nil {
			conn.Close()
			dstServer.Close()
		}
	}()
	encodeCopy(conn, dstServer, enpwd)
}
func main() {
	decodePassword := core.Password{225, 96, 201, 157, 48, 11, 5, 142, 177, 195, 163, 125, 13, 147, 164, 174, 114, 221, 243, 180, 58, 17, 224, 134, 9, 220, 138, 118, 102, 170, 169, 23, 94, 202, 240, 8, 254, 219, 146, 188, 176, 222, 44, 227, 20, 90, 145, 2, 204, 0, 250, 76, 179, 162, 140, 119, 156, 211, 245, 139, 238, 223, 3, 235, 203, 81, 130, 100, 208, 104, 25, 193, 53, 84, 77, 62, 131, 117, 63, 98, 161, 168, 173, 107, 214, 200, 38, 194, 228, 253, 7, 186, 95, 1, 207, 86, 83, 135, 115, 175, 40, 165, 217, 185, 99, 178, 21, 6, 18, 237, 97, 159, 27, 24, 45, 206, 46, 59, 43, 73, 196, 42, 184, 141, 148, 158, 50, 229, 189, 68, 89, 215, 87, 39, 121, 255, 49, 88, 160, 209, 127, 108, 72, 187, 16, 144, 80, 128, 212, 111, 137, 216, 47, 74, 67, 242, 116, 126, 199, 246, 92, 183, 239, 30, 154, 19, 4, 70, 82, 181, 236, 244, 56, 249, 233, 54, 232, 101, 136, 226, 248, 182, 166, 85, 103, 57, 79, 71, 149, 155, 35, 28, 120, 106, 109, 60, 36, 66, 241, 151, 91, 10, 12, 122, 124, 14, 143, 234, 231, 123, 26, 191, 78, 32, 31, 15, 112, 22, 133, 192, 75, 55, 150, 247, 93, 210, 171, 61, 34, 152, 129, 190, 113, 153, 33, 51, 205, 52, 167, 213, 197, 230, 198, 132, 29, 251, 41, 105, 252, 218, 64, 65, 69, 37, 110, 172}

	encodePassword := core.Password{49, 93, 47, 62, 166, 6, 107, 90, 35, 24, 201, 5, 202, 12, 205, 215, 144, 21, 108, 165, 44, 106, 217, 31, 113, 70, 210, 112, 191, 244, 163, 214, 213, 234, 228, 190, 196, 253, 86, 133, 100, 246, 121, 118, 42, 114, 116, 152, 4, 136, 126, 235, 237, 72, 175, 221, 172, 185, 20, 117, 195, 227, 75, 78, 250, 251, 197, 154, 129, 252, 167, 187, 142, 119, 153, 220, 51, 74, 212, 186, 146, 65, 168, 96, 73, 183, 95, 132, 137, 130, 45, 200, 160, 224, 32, 92, 1, 110, 79, 104, 67, 177, 28, 184, 69, 247, 193, 83, 141, 194, 254, 149, 216, 232, 16, 98, 156, 77, 27, 55, 192, 134, 203, 209, 204, 11, 157, 140, 147, 230, 66, 76, 243, 218, 23, 97, 178, 150, 26, 59, 54, 123, 7, 206, 145, 46, 38, 13, 124, 188, 222, 199, 229, 233, 164, 189, 56, 3, 125, 111, 138, 80, 53, 10, 14, 101, 182, 238, 81, 30, 29, 226, 255, 82, 15, 99, 40, 8, 105, 52, 19, 169, 181, 161, 122, 103, 91, 143, 39, 128, 231, 211, 219, 71, 87, 9, 120, 240, 242, 158, 85, 2, 33, 64, 48, 236, 115, 94, 68, 139, 225, 57, 148, 239, 84, 131, 151, 102, 249, 37, 25, 17, 41, 61, 22, 0, 179, 43, 88, 127, 241, 208, 176, 174, 207, 63, 170, 109, 60, 162, 34, 198, 155, 18, 171, 58, 159, 223, 180, 173, 50, 245, 248, 89, 36, 135}
	listen, err := net.Listen("tcp", ":8080")

	if err != nil {
		return
	}

	fmt.Println("Server is running on 8080")

	for {
		conn, _ := listen.Accept()
		go process(conn, decodePassword, encodePassword)
	}

}