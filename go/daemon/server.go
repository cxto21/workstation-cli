package daemon

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mr-kelly/mato/go/protocol"
)

type Server struct {
	SocketPath string
	Version    string

	ln   net.Listener
	wg   sync.WaitGroup
	once sync.Once
}

func DefaultSocketPath() (string, error) {
	if v := os.Getenv("MATO_SOCKET"); v != "" {
		return v, nil
	}
	if v := os.Getenv("WORKSTATION_HOME"); v != "" {
		return filepath.Join(v, "state", "daemon.sock"), nil
	}
	home, err := os.UserHomeDir()
	if err == nil {
		newPath := filepath.Join(home, ".workstation", "state", "daemon.sock")
		if _, statErr := os.Stat(newPath); statErr == nil {
			return newPath, nil
		}
	}
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		stateHome = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(stateHome, "mato", "daemon.sock"), nil
}

func (s *Server) Run(ctx context.Context) error {
	if s.SocketPath == "" {
		return errors.New("socket path is required")
	}
	if s.Version == "" {
		s.Version = "go-dev"
	}

	if err := os.MkdirAll(filepath.Dir(s.SocketPath), 0o755); err != nil {
		return fmt.Errorf("create socket dir: %w", err)
	}
	_ = os.Remove(s.SocketPath)

	ln, err := net.Listen("unix", s.SocketPath)
	if err != nil {
		return fmt.Errorf("listen unix socket: %w", err)
	}
	s.ln = ln
	_ = os.Chmod(s.SocketPath, 0o600)

	go func() {
		<-ctx.Done()
		s.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) || ctx.Err() != nil {
				break
			}
			continue
		}
		s.wg.Add(1)
		go func(c net.Conn) {
			defer s.wg.Done()
			defer c.Close()
			s.handleConn(c)
		}(conn)
	}

	s.wg.Wait()
	_ = os.Remove(s.SocketPath)
	return nil
}

func (s *Server) Close() {
	s.once.Do(func() {
		if s.ln != nil {
			_ = s.ln.Close()
		}
	})
}

func (s *Server) handleConn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var req protocol.ClientMsg
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			_ = writeServerMsg(writer, protocol.ServerMsg{
				Variant: protocol.ServerError,
				Payload: protocol.ServerErrorPayload{Message: "invalid request"},
			})
			continue
		}

		resp := s.dispatch(req)
		if err := writeServerMsg(writer, resp); err != nil {
			return
		}
	}
}

func (s *Server) dispatch(req protocol.ClientMsg) protocol.ServerMsg {
	switch req.Variant {
	case protocol.ClientHello:
		return protocol.ServerMsg{
			Variant: protocol.ServerWelcome,
			Payload: protocol.ServerWelcomePayload{Version: s.Version},
		}
	default:
		return protocol.ServerMsg{
			Variant: protocol.ServerError,
			Payload: protocol.ServerErrorPayload{Message: "unsupported in go bootstrap daemon"},
		}
	}
}

func writeServerMsg(w *bufio.Writer, msg protocol.ServerMsg) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	if err := w.WriteByte('\n'); err != nil {
		return err
	}
	return w.Flush()
}
