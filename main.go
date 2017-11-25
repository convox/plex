package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/convox/plex/client"
	"github.com/convox/plex/server"
	"github.com/convox/plex/util"
	"github.com/inconshreveable/muxado"
)

var flagLocal stringSlice
var flagRemote stringSlice

func main() {
	if err := run(); err != nil {
		stack := debug.Stack()
		fmt.Println(string(stack))
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage:\n  plex [-l local:remote]... [-r remote:local]... <command>\n  plex server\n")
		os.Exit(1)
	}

	if os.Args[1] == "server" {
		return server.New(muxado.Server(util.ReadWriteCloser{Reader: os.Stdin, WriteCloser: os.Stdout}, nil)).Run()
	}

	flag.Var(&flagLocal, "l", "forward from local to remote")
	flag.Var(&flagRemote, "r", "foward from remote to local")
	flag.Parse()

	fmt.Printf("flag.Args() = %+v\n", flag.Args())

	a, b := net.Pipe()

	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)

	cmd.Stdin = a
	cmd.Stdout = a
	cmd.Stderr = os.Stdout

	c := client.New(muxado.Client(b, nil))

	if err := cmd.Start(); err != nil {
		return err
	}

	for _, l := range flagLocal {
		parts := strings.Split(l, ":")

		if len(parts) != 2 {
			return fmt.Errorf("invalid local forward: %s", l)
		}

		if err := c.ForwardLocal(parts[0], parts[1]); err != nil {
			return err
		}
	}

	for _, l := range flagRemote {
		parts := strings.Split(l, ":")

		if len(parts) != 2 {
			return fmt.Errorf("invalid remote forward: %s", l)
		}

		if err := c.ForwardRemote(parts[0], parts[1]); err != nil {
			return err
		}
	}

	s := server.New(c.Session)

	return s.Run()
}

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}
