package main

import (
	"fmt"
	"os"

	"github.com/zaydek/go-ssr/cmd/export"
	"github.com/zaydek/go-ssr/cmd/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, `
	Usage:

		go run main.go server  Run the SSR server
		go run main.go export  Export the SSR server

`)
		return
	}
	switch os.Args[1] {
	case "export":
		export.Run()
	case "server":
		server.Run()
	default:
		panic("Internal error")
	}
}
