package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
)

// TCPServer simple tcp server for commands
func TCPServer() (ret int) {
	ln, err := net.Listen("tcp", ":"+TCPPort)
	if err != nil {
		log.Println("[control] error listening:", err.Error())
		return 1
	}
	log.Printf("[control] listening to port %s for connections...", TCPPort)
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		bot := NewBot()
		bot.inconn = conn
		if err != nil {
			fmt.Println("[control] error accepting: ", err.Error())
			return 1
		}
		go handleRequest(conn, bot)
	}
}

func handleRequest(conn net.Conn, bot *Bot) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()

		if err != nil {
			fmt.Println("[control] read error:", err)
			fmt.Fprintf(conn, "[control] read error: %s", err)
			conn.Close()
			return
		}
		log.Println("[control] " + line)
		handleMessage(line, bot)
	}
}

// Handle an IRC received from a bot
func handleMessage(message string, bot *Bot) {
	remoteAddr := bot.inconn.RemoteAddr().String()
	remoteAddrIP := strings.Split(remoteAddr, ":")

	if strings.Contains(message, "JOIN ") {
		joinComm := strings.Split(message, "JOIN ")
		go bot.join(joinComm[1])
	} else if strings.Contains(message, "PASS ") {
		passComm := strings.Split(message, "PASS ")
		passwordParts := strings.Split(passComm[1], ";")
		if passwordParts[0] == TCPPass {
			bot.oauth = passwordParts[1]
			log.Printf("[control] authenticated! %s\n", remoteAddrIP)
		} else {
			log.Printf("[control] invalid broker pass! %s\n", remoteAddrIP)
			bot.inconn.Close()
			return
		}
	} else if strings.Contains(message, "NICK ") {
		nickComm := strings.Split(message, "NICK ")
		bot.nick = nickComm[1]

		if bot.oauth != "" {
			bot.CreateConnection()
			go bot.CreateGroupConnection()
		}
	} else if strings.Contains(message, "USER ") {
		if bot.nick != "" {
			nickComm := strings.Split(message, "USER ")
			bot.nick = nickComm[1]

			if bot.oauth != "" {
				bot.CreateConnection()
				go bot.CreateGroupConnection()
			}
		}
	} else if strings.Contains(message, "PRIVMSG #jtv :/w ") {
		privmsgComm := strings.Split(message, "PRIVMSG #jtv :")
		go bot.Whisper(privmsgComm[1])
	} else {
		bot.Message(message)
	}
}
