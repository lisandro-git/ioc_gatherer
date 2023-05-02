package server

import (
	"fmt"
	"net"
	"os"
	"main/cmd"
)

func StartLocalServer() net.Listener {
	listener, err := net.Listen(cmd.TYPE, cmd.REMOTEHOST+cmd.REMORTPORT)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on " + cmd.REMOTEHOST + cmd.REMORTPORT)
	return listener
}

func ReadData(conn net.Conn) {
	var buffer [cmd.BUFFERSIZE]byte
	for {
		read, err := conn.Read(buffer[:])
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed")
				conn.Close()
				return
			} else {
				panic(err)
			}
		}
		fmt.Println("Received Data: ", string(buffer[:read]))
	}
}
