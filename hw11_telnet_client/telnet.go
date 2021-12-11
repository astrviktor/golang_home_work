package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"sync"
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
	mu      *sync.Mutex
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
	return FromInToOut(stc.mu, stc.in, stc.conn)
}

func (stc *SimpleTelnetClient) Receive() error {
	return FromInToOut(stc.mu, stc.conn, stc.out)
}

func FromInToOut(mu *sync.Mutex, in io.ReadCloser, out io.Writer) error {
	scanner := bufio.NewScanner(in)

	mu.Lock()
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
	mu.Unlock()

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	mu := sync.Mutex{}
	client := SimpleTelnetClient{nil, address, timeout, in, out, &mu}
	return &client
}
