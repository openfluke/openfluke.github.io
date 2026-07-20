package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/donate"
)

func main() {
	fmt.Println("default port", donate.DefaultPort)
	// ln, err := donate.ServeTCP(donate.ServerOptions{Addr: ":17001", Mode: donate.ServerLocalLM})
}
