package daemon

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/mr-kelly/mato/go/protocol"
)

func TestHelloHandshake(t *testing.T) {
	socket := filepath.Join(t.TempDir(), "daemon.sock")
	srv := &Server{SocketPath: socket, Version: "0.9.6-go"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() { errCh <- srv.Run(ctx) }()

	deadline := time.Now().Add(2 * time.Second)
	for {
		if time.Now().After(deadline) {
			t.Fatal("server did not start listening in time")
		}
		conn, err := net.Dial("unix", socket)
		if err == nil {
			_ = conn.Close()
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	conn, err := net.Dial("unix", socket)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}

	req := protocol.ClientMsg{
		Variant: protocol.ClientHello,
		Payload: protocol.ClientHelloPayload{Version: "test-client"},
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request failed: %v", err)
	}
	if _, err := conn.Write(append(data, '\n')); err != nil {
		t.Fatalf("write request failed: %v", err)
	}

	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatalf("read response failed: %v", err)
	}

	var resp protocol.ServerMsg
	if err := json.Unmarshal([]byte(line), &resp); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if resp.Variant != protocol.ServerWelcome {
		t.Fatalf("unexpected variant: %s", resp.Variant)
	}
	payload, ok := resp.Payload.(protocol.ServerWelcomePayload)
	if !ok {
		t.Fatalf("unexpected payload type: %T", resp.Payload)
	}
	if payload.Version != "0.9.6-go" {
		t.Fatalf("unexpected version: %s", payload.Version)
	}

	_ = conn.Close()
	cancel()
	if err := <-errCh; err != nil {
		t.Fatalf("server ended with error: %v", err)
	}
}
