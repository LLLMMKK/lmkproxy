package core

import (
	"fmt"
	"math/rand"
)

const PasswordLength = 256

type Password [PasswordLength]byte

func Init() Password {
	a := rand.Perm(PasswordLength)
	password := Password{}
	for i := 0; i < PasswordLength; i++ {
		password[i] = byte(a[i])
	}
	fmt.Println(password)
	return password
}

func Decode(depwd Password, buf []byte) {
	for i, v := range buf {
		buf[i] = depwd[v]
	}
}

func Encode(enpwd Password, buf []byte) {
	for i, v := range buf {
		buf[i] = enpwd[v]
	}
}
