package liquidsoap

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"go.uber.org/zap"
)

type Telnet struct {
	addr   string
	logger *zap.Logger
	mutex  sync.Mutex
}

func NewTelnet(addr string, logger *zap.Logger) *Telnet {
	return &Telnet{
		addr:   addr,
		logger: logger.Named("liq"),
	}
}

func (t *Telnet) Command(command string) (string, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.logger.Debug("connecting using telnet", zap.String("addr", t.addr))
	conn, err := net.Dial("tcp", t.addr)
	if err != nil {
		return "", fmt.Errorf("failed to connect to liquidsoap telnet server: %w", err)
	}
	defer conn.Close()

	t.logger.Debug("sending command", zap.String("command", command))
	fmt.Fprintf(conn, "%s\n", command)

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("telnet liq error: %w", err)
	}
	t.logger.Debug("received", zap.String("message", message))

	return message, nil
}
