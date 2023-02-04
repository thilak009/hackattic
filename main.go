package main

import (
	"fmt"
	"os"

	"github.com/thilak009/hackattic/readqr"
	"github.com/thilak009/hackattic/unpack"
)

func main() {
	token := os.Getenv("token")
	if token == "" {
		fmt.Println("token is not present")
		return
	}
	unpack.Unpack(token)
	readqr.ReadQR(token)
}
