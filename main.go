package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strings"
)

type Header map[string]string

const (
	Host          = "Host"
	ContentType   = "Content-Type"
	ContentLength = "Content-Length"
	UserAgent     = "User-Agent"
	Accept        = "Accept"
	Httpv1        = "HTTP/1.1"
	Get           = "GET"
	Post          = "POST"
	Delete        = "DELETE"
	Put           = "PUT"
	Patch         = "PATCH"
)

var AllowedMethods []string = []string{Get, Post, Delete, Patch}

type Request struct {
	method  string
	path    string
	version string
	header  Header
	body    []byte
}

func main() {
	ln, err := net.Listen("tcp", ":6703")
	if err != nil {
		slog.Error("listen", "error", err)
		return
	}
	slog.Info("Listening on port 6703")

	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("accept", "error", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	header := make(Header)
	reader := bufio.NewReader(conn)
	metadata, _ := reader.ReadString('\n')
	mTokens := strings.SplitN(metadata, " ", 3)
	// method := strings.TrimSpace(mTokens[0])
	// path := strings.TrimSpace(mTokens[1])
	if version := strings.TrimSpace(mTokens[2]); version != Httpv1 {
		slog.Info("HTTP Version not supported, please revert to 1.0")
		conn.Close()
		return
	}

	for {
		line, _ := reader.ReadString('\n')
		if line == "\r\n" {
			break
		}
		tokens := strings.SplitN(line, ":", 2)
		header[tokens[0]] = strings.TrimSpace(tokens[1])
	}

	conn.Close()
	return
}
