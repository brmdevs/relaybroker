package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
)

// Bot struct for main config
type Bot struct {
	server      string
	groupserver string
	port        string
	connlist    []net.Conn
}

// NewBot main config
func NewBot() *Bot {
	return &Bot{
		server:      "irc.twitch.tv",
		groupserver: "group.tmi.twitch.tv",
		port:        "6667",
		connlist:    make([]net.Conn, 0),
	}
}

// CreateConnection Add a new connection
func (bot *Bot) CreateConnection() {
	conn, err := net.Dial("tcp", bot.server+":"+bot.port)
	if err != nil {
		log.Fatal("unable to connect to IRC server ", err)
	}
	log.Printf("Connected to IRC server %s (%s)\n", bot.server, conn.RemoteAddr())
	bot.connlist = append(bot.connlist, conn)

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()
		if err != nil {
			break // break loop on errors
		}
		if strings.Contains(line, "PING") {
			pongdata := strings.Split(line, "PING ")
			fmt.Fprintf(conn, "PONG %s\r\n", pongdata[1])
		}
		bot.Handle(line)
	}
}

func main() {
	ircbot := NewBot()
	TCPServer(ircbot)
}

// HandleJoin will slowly join all channels given
// 45 per 11 seconds to deal with twitch ratelimits
func (bot *Bot) HandleJoin(channels []string) {
	for _, channel := range channels {
		for _, conn := range bot.connlist {
			fmt.Println("Joining " + channel)
			fmt.Println(conn)
			fmt.Fprintf(conn, "JOIN %s\r\n", channel)
		}
	}
}

// WriteToAllConns writes message to all connections for now
func (bot *Bot) WriteToAllConns(message string) {
	for _, conn := range bot.connlist {
		fmt.Fprintf(conn, message+"\r\n")
	}
}

// Message to send a message
func (bot *Bot) Message(channel string, message string) {
	if message == "" {
		return
	}
	fmt.Printf("Bot: " + message + "\n")

	/* Find a suitable connection to use */
	for _, conn := range bot.connlist {
		fmt.Fprintf(conn, "PRIVMSG %s :%s\r\n", channel, message)
	}
}

// Handle handles messages from irc
func (bot *Bot) Handle(line string) {
	if strings.Contains(line, ".tmi.twitch.tv PRIVMSG ") {
		messageTMISplit := strings.Split(line, ".tmi.twitch.tv PRIVMSG ")
		messageChannelRaw := strings.Split(messageTMISplit[1], " :")
		channel := messageChannelRaw[0]
		go bot.ProcessMessage(channel, line)
	} else if strings.Contains(line, ":tmi.twitch.tv ROOMSTATE") {
		messageTMISplit := strings.Split(line, ":tmi.twitch.tv ROOMSTATE ")
		channel := messageTMISplit[1]
		go bot.ProcessMessage(channel, line)
	}
}

// ProcessMessage push message to local irc chat
func (bot *Bot) ProcessMessage(channel string, message string) {
	fmt.Println(channel + " ::: " + message)
}
