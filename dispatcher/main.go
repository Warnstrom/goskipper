package main

import (
	"fmt"
	"net"

	hook "github.com/robotn/gohook"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "jonatan.net:9009")
	if err != nil {
		fmt.Print("Couldn't resolve the address")
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Print("Couldn't connect to the server")
	}
	for {
		hook.Register(hook.KeyDown, []string{"shift", "space"}, func(e hook.Event) {

			fmt.Fprintf(conn, "text"+"\n")
			hook.End()
		})

		s := hook.Start()
		<-hook.Process(s)
	}
}
