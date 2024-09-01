package request

import (
	"bufio"
	"net"
	"strings"
	"log/slog"
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
	request := Request{}
	reader := bufio.NewReader(conn)
	metadata, err := reader.ReadString('\n')
	if err != nil {
		slog.Info("An error occured while reading the header", "error", err)
		return Request{}, err
	}
	mtokens := strings.SplitN(metadata, " ", 3)
	request.Method = strings.TrimSpace(mtokens[0])
	request.Path = strings.TrimSpace(mtokens[1])
	request.Version = strings.TrimSpace(mtokens[2])

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
	request.Header = header
	return request, nil
}
