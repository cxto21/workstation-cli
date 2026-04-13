package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mr-kelly/mato/go/daemon"
	"github.com/mr-kelly/mato/go/protocol"
)

func main() {
	mode := flag.String("mode", "hello", "Mode: hello | daemon")
	socket := flag.String("socket", "", "Unix socket path (optional)")
	version := flag.String("version", "go-dev", "Version string used by daemon welcome")
	flag.Parse()

	path := *socket
	if path == "" {
		resolved, err := daemon.DefaultSocketPath()
		if err != nil {
			fatal(err)
		}
		path = resolved
	}

	switch strings.ToLower(*mode) {
	case "daemon":
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		srv := &daemon.Server{SocketPath: path, Version: *version}
		fmt.Printf("mato-go daemon listening on %s\n", path)
		if err := srv.Run(ctx); err != nil {
			fatal(err)
		}
	case "hello":
		if err := doHello(path, *version); err != nil {
			fatal(err)
		}
	default:
		fatal(fmt.Errorf("unsupported mode: %s", *mode))
	}
}

func doHello(socketPath, clientVersion string) error {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return fmt.Errorf("dial daemon: %w", err)
	}
	defer conn.Close()

	req := protocol.ClientMsg{
		Variant: protocol.ClientHello,
		Payload: protocol.ClientHelloPayload{Version: clientVersion},
	}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	if _, err := conn.Write(append(data, '\n')); err != nil {
		return err
	}

	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return err
	}

	var resp protocol.ServerMsg
	if err := json.Unmarshal([]byte(line), &resp); err != nil {
		return err
	}
	if resp.Variant != protocol.ServerWelcome {
		return errors.New("unexpected response from daemon")
	}
	welcome, ok := resp.Payload.(protocol.ServerWelcomePayload)
	if !ok {
		return errors.New("unexpected welcome payload")
	}

	fmt.Printf("mato-go hello ok: daemon_version=%s socket=%s\n", welcome.Version, socketPath)
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
