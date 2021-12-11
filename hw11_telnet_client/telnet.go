package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type SimpleTelnetClient struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (stc *SimpleTelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", stc.address, stc.timeout)
	if err != nil {
		return err
	}

	stc.conn = conn
	return nil
}

func (stc *SimpleTelnetClient) Close() error {
	return stc.conn.Close()
}

func (stc *SimpleTelnetClient) Send() error {
	return FromInToOut(stc.in, stc.conn)
}

func (stc *SimpleTelnetClient) Receive() error {
	return FromInToOut(stc.conn, stc.out)
}

func FromInToOut(in io.ReadCloser, out io.Writer) error {
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		if errors.Is(scanner.Err(), io.EOF) {
			_, err := out.Write([]byte("^D\n...EOF"))
			if err != nil {
				return err
			}
			return io.EOF
		}

		str := scanner.Text() + "\n"
		_, err := out.Write([]byte(str))
		if err != nil {
			return err
		}
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	client := SimpleTelnetClient{nil, address, timeout, in, out}
	return &client
}
