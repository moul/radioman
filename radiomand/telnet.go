package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
)

type LiquidsoapTelnet struct {
	Host string
	Port int
	Conn net.Conn
}

func NewLiquidsoapTelnet(host string, port int) *LiquidsoapTelnet {
	return &LiquidsoapTelnet{
		Host: host,
		Port: port,
	}
}

func (t *LiquidsoapTelnet) Dest() string {
	return fmt.Sprintf("%s:%d", t.Host, t.Port)
}

func (t *LiquidsoapTelnet) Open() error {
	logrus.Infof("Connecting to Liquidsoap telnet: %s:%d", t.Host, t.Port)
	conn, err := net.Dial("tcp", t.Dest())
	if err != nil {
		return err
	}
	t.Conn = conn
	return nil
}

func (t *LiquidsoapTelnet) Close() {
	if t.Conn != nil {
		t.Conn.Close()
	}
}

func (t *LiquidsoapTelnet) Command(command string) (string, error) {
	if err := t.Open(); err != nil {
		logrus.Errorf("Failed to connect to liquidsoap telnet socket")
		return "", err
	}
	defer t.Close()
	logrus.Infof("Sending to telnet: %q", command)
	fmt.Fprintf(t.Conn, "%s\n", command)
	message, err := bufio.NewReader(t.Conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	fmt.Printf("Received from telnet: %q", message)
	return message, nil
}
