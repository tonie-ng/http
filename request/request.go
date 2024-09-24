package request

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
)

type (
	Header  map[string]string
	Request struct {
		Method  string
		Path    string
		Version string
		Header  Header
		Body    []byte
	}
)

func ParseRequest(conn net.Conn) (Request, error) {
	req := Request{}
	reader := bufio.NewReader(conn)
	metadata, err := reader.ReadString('\n')
	if err != nil {
		slog.Info("An error occured while reading the header", "error", err)
		return Request{}, err
	}
	mtokens := strings.SplitN(metadata, " ", 3)
	req.Method = strings.TrimSpace(mtokens[0])
	req.Path = strings.TrimSpace(mtokens[1])
	req.Version = strings.TrimSpace(mtokens[2])

	header := make(Header)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			slog.Info("An error occured while reading the header", "error:", err)
			return Request{}, err
		}
		if line == "\r\n" {
			break
		}
		tokens := strings.SplitN(line, ":", 2)
		header[tokens[0]] = strings.TrimSpace(tokens[1])
	}
	req.Header = header

	// in future we might have to support PATCH or PUT
	if req.Method == "POST" {
		cLen, err := strconv.Atoi(header["Content-Length"])
		if err != nil || cLen < 1 {
			err = fmt.Errorf("No or invalid content length on %s request", req.Method)
			return Request{}, err
		}

		body := make([]byte, cLen)
		reader.Read(body)
		req.Body = body
	}
	return req, nil
}
