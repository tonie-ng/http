package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"
)

type Header map[string]string

const (
	Host                = "Host"
	ContentType         = "Content-Type"
	ContentLength       = "Content-Length"
	UserAgent           = "User-Agent"
	Accept              = "Accept"
	NotFound            = "404 Not Found"
	BadRequest          = "400 Bad Request"
	NotAllowed          = "405 Method Not Allowed"
	InternalServerError = "500 Internal Server Error"
	Ok                  = "200 OK"
	Get                 = "GET"
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
	defer conn.Close()

	header := make(Header)
	reader := bufio.NewReader(conn)
	metadata, err := reader.ReadString('\n')
	if err != nil {
		slog.Info("An error occured while reading the header", "error:", err)
		WriteHeader(conn, "", BadRequest, 0)
		return
	}
	mTokens := strings.SplitN(metadata, " ", 3)
	if method := strings.TrimSpace(mTokens[0]); method != "GET" {
		slog.Info("Method not supported (at least for now)")
		WriteHeader(conn, "", BadRequest, 0)
		return
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			slog.Info("An error occured while reading the header", "error:", err)
			WriteHeader(conn, "", BadRequest, 0)
			return
		}
		if line == "\r\n" {
			break
		}
		tokens := strings.SplitN(line, ":", 2)
		header[tokens[0]] = strings.TrimSpace(tokens[1])
	}

	// this path should be sanitized to avoid malicious attacks
	reqFilePath := strings.TrimSpace(mTokens[1])
	fileInfo, filePath, err := findFile(reqFilePath)
	if err != nil {
		slog.Info("An error occured opening the file", "error", err)
		WriteHeader(conn, "", NotFound, 0)
		return
	}

	data, _ := os.ReadFile(filePath)
	WriteHeader(conn, fileInfo.Name(), Ok, fileInfo.Size())
	conn.Write(data)
	return
}

func WriteHeader(conn net.Conn, filename, status string, contentLength int64) error {
	res := fmt.Sprintf("HTTP/1.1 %s \r\nContent-Type: %s\r\nContent-Length: %d\r\nDate: %s\r\n\r\n", status, GetContentType(filename), contentLength, time.Now().Format(time.RFC1123))
	conn.Write([]byte(res))
	return nil
}

func GetContentType(filename string) string {
	ext := strings.Split(filename, ".")
	switch ext[len(ext)-1] {
	case "js":
		return "application/javascript"
	case "jpg":
		return "image/jpg"
	case "png":
		return "image/png"
	case "html":
		return "text/html"
	default:
		return "text/plain"
	}
}

func findFile(filename string) (os.FileInfo, string, error) {
	if filename == "/" {
		filename = "index.html"
	}
	if filename[0] == '/' {
		filename = filename[1:]
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, "", err
	}

	if fileInfo.IsDir() {
		if filename[len(filename)-1] == '/' {
			filename = fmt.Sprintf("%sindex.html", filename)
		} else {
			filename = fmt.Sprintf("%s/index.html", filename)
		}
		fileInfo, err = os.Stat(filename)
		if err != nil {
			return nil, "", err
		}
	}
	return fileInfo, filename, nil
}
