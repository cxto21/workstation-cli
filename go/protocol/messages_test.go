package protocol

import (
	"encoding/json"
	"testing"
)

func TestClientSpawnRoundtrip(t *testing.T) {
	msg := ClientMsg{
		Variant: ClientSpawn,
		Payload: ClientSpawnPayload{
			TabID: "test123",
			Rows:  24,
			Cols:  80,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var back ClientMsg
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if back.Variant != ClientSpawn {
		t.Fatalf("unexpected variant: %s", back.Variant)
	}

	payload, ok := back.Payload.(ClientSpawnPayload)
	if !ok {
		t.Fatalf("unexpected payload type: %T", back.Payload)
	}

	if payload.TabID != "test123" || payload.Rows != 24 || payload.Cols != 80 {
		t.Fatalf("unexpected payload content: %+v", payload)
	}
}

func TestClientInputUsesByteArrayNumbers(t *testing.T) {
	msg := ClientMsg{
		Variant: ClientInput,
		Payload: ClientInputPayload{
			TabID: "abc",
			Data:  ByteArray{65, 66, 67},
		},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var obj map[string]map[string]any
	if err := json.Unmarshal(data, &obj); err != nil {
		t.Fatalf("unmarshal object failed: %v", err)
	}

	input := obj["Input"]
	rawData, ok := input["data"].([]any)
	if !ok {
		t.Fatalf("data should be numeric array, got %T", input["data"])
	}
	if len(rawData) != 3 || rawData[0].(float64) != 65 || rawData[1].(float64) != 66 || rawData[2].(float64) != 67 {
		t.Fatalf("unexpected data values: %+v", rawData)
	}
}

func TestServerWelcomeRoundtrip(t *testing.T) {
	msg := ServerMsg{
		Variant: ServerWelcome,
		Payload: ServerWelcomePayload{Version: "0.1.0"},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var back ServerMsg
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if back.Variant != ServerWelcome {
		t.Fatalf("unexpected variant: %s", back.Variant)
	}
	payload, ok := back.Payload.(ServerWelcomePayload)
	if !ok {
		t.Fatalf("unexpected payload type: %T", back.Payload)
	}
	if payload.Version != "0.1.0" {
		t.Fatalf("unexpected version: %s", payload.Version)
	}
}

func TestInvalidEnvelopeRejected(t *testing.T) {
	input := []byte(`{"Hello":{"version":"1.0"},"Spawn":{"tab_id":"x","rows":1,"cols":1}}`)
	var msg ClientMsg
	if err := json.Unmarshal(input, &msg); err == nil {
		t.Fatal("expected invalid envelope error")
	}
}

func TestClientScrollRoundtrip(t *testing.T) {
	msg := ClientMsg{
		Variant: ClientScroll,
		Payload: ClientScrollPayload{TabID: "t1", Delta: -5},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var back ClientMsg
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if back.Variant != ClientScroll {
		t.Fatalf("unexpected variant: %s", back.Variant)
	}
	p, ok := back.Payload.(ClientScrollPayload)
	if !ok {
		t.Fatalf("unexpected payload type: %T", back.Payload)
	}
	if p.TabID != "t1" || p.Delta != -5 {
		t.Fatalf("unexpected payload: %+v", p)
	}
}

func TestClientGetInputModesRoundtrip(t *testing.T) {
	msg := ClientMsg{
		Variant: ClientGetInputModes,
		Payload: ClientGetInputModesPayload{TabID: "tab-1"},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var back ClientMsg
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if back.Variant != ClientGetInputModes {
		t.Fatalf("unexpected variant: %s", back.Variant)
	}
	p, ok := back.Payload.(ClientGetInputModesPayload)
	if !ok {
		t.Fatalf("unexpected payload type: %T", back.Payload)
	}
	if p.TabID != "tab-1" {
		t.Fatalf("unexpected tab_id: %s", p.TabID)
	}
}

func TestServerProcessStatusRoundtrip(t *testing.T) {
	msg := ServerMsg{
		Variant: ServerProcessStatus,
		Payload: ServerProcessStatusPayload{Tabs: [][2]any{{"tab-a", 1234.0}}},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var back ServerMsg
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if back.Variant != ServerProcessStatus {
		t.Fatalf("unexpected variant: %s", back.Variant)
	}
	p, ok := back.Payload.(ServerProcessStatusPayload)
	if !ok {
		t.Fatalf("unexpected payload type: %T", back.Payload)
	}
	if len(p.Tabs) != 1 || p.Tabs[0][0] != "tab-a" {
		t.Fatalf("unexpected tabs payload: %+v", p.Tabs)
	}
}
