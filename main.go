package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
)

type Header map[string]string

const (
	Host          = "Host"
	ContentType   = "Content-Type"
	ContentLength = "Content-Length"
	UserAgent     = "User-Agent"
	Accept        = "Accept"
	Get           = "GET"
)

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
	metadata, err := reader.ReadString('\n')
	if err != nil {
		slog.Info("An error occured while reading the header", "error:", err)
		conn.Close()
		return
	}
	mTokens := strings.SplitN(metadata, " ", 3)
	if method := strings.TrimSpace(mTokens[0]); method != "GET" {
		slog.Info("Method not supported (at least for now)")
		conn.Close()
		return
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			slog.Info("An error occured while reading the header", "error:", err)
			conn.Close()
			return
		}
		if line == "\r\n" {
			break
		}
		tokens := strings.SplitN(line, ":", 2)
		header[tokens[0]] = strings.TrimSpace(tokens[1])
	}

	// this path should be sanitized to avoid malicious attacks
	filepath := strings.TrimSpace(mTokens[1])
	fileInfo, err := findFile(filepath)
	if err != nil {
		slog.Info("An error occured opening the file", "error", err)
		conn.Close()
		return
	}

	os.ReadFile(fileInfo.Name())
	conn.Close()
	return
}


func findFile(filename string) (os.FileInfo, error) {
	if filename == "/" {
		filename = "index.html"
	}
	if filename[0] == '/' {
		filename = filename[1:]
	}
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		if filename[len(filename)-1] == '/' {
			filename = fmt.Sprintf("%sindex.html", filename)
		} else {
			filename = fmt.Sprintf("%s/index.html", filename)
		}
		fileInfo, err = os.Stat(filename)
		if err != nil {
			return nil, err
		}
	}
	return fileInfo, err
}
