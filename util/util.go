package util

import "io"

type ReadWriteCloser struct {
	io.Reader
	io.WriteCloser
}

func Pipe(a, b io.ReadWriter) error {
	ch := make(chan error)

	go halfPipe(a, b, ch)
	go halfPipe(b, a, ch)

	if err := <-ch; err != nil {
		return err
	}

	if err := <-ch; err != nil {
		return err
	}

	return nil
}

func halfPipe(w io.Writer, r io.Reader, ch chan error) {
	_, err := io.Copy(w, r)
	ch <- err
}
