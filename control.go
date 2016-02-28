package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	auth = make([]string, 8)
)

// TCPServer simple tcp server for commands
func TCPServer(ircbot *Bot) {
	ln, err := net.Listen("tcp", ":"+TCPPort)
	log.Println("TCP Server listening on port", TCPPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		}
		go handleRequest(conn, ircbot)
	}

}

func handleRequest(conn net.Conn, ircbot *Bot) {

	message, _ := bufio.NewReader(conn).ReadString('\n')
	remoteAddr := conn.RemoteAddr().String()
	remoteAddrIP := strings.Split(remoteAddr, ":")
	fmt.Println(message)

	if stringInSlice(remoteAddrIP[0], auth) {
		handleMessage(message, ircbot)
		conn.Write([]byte("Message received"))
	} else if message == "AUTH "+TCPPass {
		auth = append(auth, remoteAddrIP[0])
		fmt.Println(auth)
		conn.Write([]byte("Authenticated\r\n"))
	} else {
		conn.Write([]byte("not authenticated use \"AUTH password\" to authenticate"))
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func handleMessage(message string, ircbot *Bot) {
	fmt.Println(message)
	ircbot.WriteToAllConns(message)
	// if strings.Contains(message, "JOIN ") {
	// 	joinComm := strings.Split(message, "JOIN ")
	// 	channels := strings.Split(joinComm[1], " ")
	// 	go ircbot.HandleJoin(channels)
	// } else if strings.Contains(message, "PRIVMSG ") {
	// 	privmsgComm := strings.Split(message, "PRIVMSG ")
	// 	remainingString := strings.Split(privmsgComm[1], " :")
	// 	channel := remainingString[0]
	// 	message := remainingString[1]
	// 	go ircbot.Message(channel, message)
	// }
}
