package server

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/convox/plex/client"
	"github.com/convox/plex/util"
	"github.com/inconshreveable/muxado"
)

type Server struct {
	muxado.Session
}

func New(s muxado.Session) *Server {
	return &Server{
		Session: s,
	}
}

func (s *Server) Run() error {
	c := client.New(s.Session)

	for {
		cn, err := s.Session.Accept()
		if err != nil {
			return err
		}

		go s.handleConnection(cn, c)
	}
}

func (s *Server) handleConnection(cn net.Conn, c *client.Client) error {
	defer cn.Close()

	header, err := readUntil(cn, '\n')
	if err != nil {
		return err
	}

	switch parts := strings.Split(string(header), " "); parts[0] {
	case "connect":
		out, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", parts[1]))
		if err != nil {
			return err
		}
		defer out.Close()

		return util.Pipe(cn, out)
	case "forward":
		go c.ForwardLocal(parts[1], parts[2])
	default:
		return fmt.Errorf("unknown header: %s", header)
	}

	return nil
}

func readUntil(r io.Reader, end byte) ([]byte, error) {
	header := []byte{}
	buf := make([]byte, 1)

	for {
		_, err := r.Read(buf)
		if err != nil {
			return nil, err
		}
		if buf[0] == end {
			return header, nil
		}

		header = append(header, buf[0])
	}
}
