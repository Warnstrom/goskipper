package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	hook "github.com/robotn/gohook"
)

var (
	tcpBaseUrl      = "jonatan.me:9009"
	keys            = []string{"shift", "tab"}         // Key(s) that it listens to
	version         = strings.Split(hook.Version, ",") // Gohook version
	keyDownListener = uint8(hook.KeyDown)              // listener event that listens to the key(s)
)

func main() {
	log.Printf("Gohook version %s is running\n", version[0])

	tcpAddr, err := net.ResolveTCPAddr("tcp", tcpBaseUrl)
	if err != nil {
		log.Fatalln("Couldn't resolve the address")
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalln("Couldn't connect to the server")
	} else {
		log.Println("Server successfully connected")
	}
	registerEvent(conn, keys)
}

func registerEvent(conn *net.TCPConn, keys []string) {
	for {
		hook.Register(keyDownListener, keys, func(e hook.Event) {
			fmt.Fprintf(conn, "text"+"\n")
			log.Printf("%s pressed: %s\n", keys[0], e.String())
			hook.End()
		})
		s := hook.Start()
		<-hook.Process(s)
	}
}
