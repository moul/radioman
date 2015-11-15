package liquidsoap

import (
	"bufio"
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
)

type Telnet struct {
	Host string
	Port int
	Conn net.Conn
}

func NewTelnet(host string, port int) *Telnet {
	return &Telnet{
		Host: host,
		Port: port,
	}
}

func (t *Telnet) Dest() string {
	return fmt.Sprintf("%s:%d", t.Host, t.Port)
}

func (t *Telnet) Open() error {
	logrus.Debugf("Connecting to Liquidsoap telnet: %s:%d", t.Host, t.Port)
	conn, err := net.Dial("tcp", t.Dest())
	if err != nil {
		return err
	}
	t.Conn = conn
	return nil
}

func (t *Telnet) Close() {
	if t.Conn != nil {
		t.Conn.Close()
	}
}

func (t *Telnet) Command(command string) (string, error) {
	logrus.Debugf("Sending to telnet: %q", command)
	fmt.Fprintf(t.Conn, "%s\n", command)
	message, err := bufio.NewReader(t.Conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	logrus.Debugf("Received from telnet: %q", message)
	return message, nil
}
