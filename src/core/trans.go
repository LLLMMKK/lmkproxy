package core

import (
	"io"
	"net"
)

func DecodeRead(conn net.Conn, depwd Password) (int, []byte, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	Decode(depwd, buf)
	return n, buf, err
}

func EncodeRead(conn net.Conn, enpwd Password) (int, []byte, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	Encode(enpwd, buf)
	return n, buf, err
}

func EncodeWrite(conn net.Conn, enpwd Password, buf []byte) (int, error) {
	Encode(enpwd, buf)
	n, err := conn.Write(buf)
	return n, err
}

func DecodeCopy(dst net.Conn, src net.Conn, depwd Password) error {
	for {
		n, buf, err := DecodeRead(src, depwd)
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

func EncodeCopy(dst net.Conn, src net.Conn, enpwd Password) error {

	for {

		n, buf, err := EncodeRead(src, enpwd)
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
