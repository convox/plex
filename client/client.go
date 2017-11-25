package client

import (
	"fmt"
	"io"
	"net"

	"github.com/inconshreveable/muxado"
)

type Client struct {
	muxado.Session
}

func New(s muxado.Session) *Client {
	return &Client{
		Session: s,
	}
}

func (c *Client) ForwardLocal(local, remote string) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", local))
	if err != nil {
		return err
	}

	go c.handleListener(ln, remote)

	return nil
}

func (c *Client) ForwardRemote(remote, local string) error {
	cn, err := c.Open()
	if err != nil {
		return err
	}

	fmt.Fprintf(cn, "forward %s %s\n", remote, local)

	return nil
}

func (c *Client) handleListener(ln net.Listener, remote string) error {
	for {
		cn, err := ln.Accept()
		if err != nil {
			return err
		}

		// fmt.Fprintf(os.Stderr, "cn = %+v\n", cn)

		go c.handleConnection(cn, remote)
	}
}

func (c *Client) handleConnection(cn net.Conn, remote string) error {
	defer cn.Close()

	st, err := c.Session.Open()
	if err != nil {
		return err
	}

	defer st.Close()

	// fmt.Fprintf(os.Stderr, "connect %v -> :%s\n", cn.LocalAddr(), remote)

	fmt.Fprintf(st, "connect %s\n", remote)

	go io.Copy(cn, st)
	io.Copy(st, cn)

	return nil
}
