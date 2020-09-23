package liquidsoap

import (
	"bufio"
	"fmt"
	"net"

	"go.uber.org/zap"
)

type Telnet struct {
	addr   string
	conn   net.Conn
	logger *zap.Logger
}

func NewTelnet(addr string, logger *zap.Logger) *Telnet {
	return &Telnet{
		addr:   addr,
		logger: logger.Named("liq"),
	}
}

func (t *Telnet) Open() error {
	t.logger.Debug("connecting using telnet", zap.String("addr", t.addr))
	var err error
	t.conn, err = net.Dial("tcp", t.addr)
	return err
}

func (t *Telnet) Close() {
	if t.conn != nil {
		t.conn.Close()
	}
}

func (t *Telnet) Command(command string) (string, error) {
	t.logger.Debug("send command", zap.String("command", command))
	fmt.Fprintf(t.conn, "%s\n", command)
	message, err := bufio.NewReader(t.conn).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("telnet liq error: %w", err)
	}
	t.logger.Debug("received", zap.String("message", message))
	return message, nil
}
