package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrInvalidEnvelope = errors.New("protocol envelope must contain exactly one variant")
	ErrUnknownVariant  = errors.New("unknown protocol variant")
)

// ByteArray preserves Rust Vec<u8> JSON shape as numeric array.
// We intentionally use []int to avoid Go's default []byte -> base64 JSON behavior.
type ByteArray []int

type ClientVariant string

const (
	ClientHello            ClientVariant = "Hello"
	ClientSpawn            ClientVariant = "Spawn"
	ClientInput            ClientVariant = "Input"
	ClientPaste            ClientVariant = "Paste"
	ClientGetInputModes    ClientVariant = "GetInputModes"
	ClientResize           ClientVariant = "Resize"
	ClientGetScreen        ClientVariant = "GetScreen"
	ClientGetProcessStatus ClientVariant = "GetProcessStatus"
	ClientSubscribe        ClientVariant = "Subscribe"
	ClientGetCwd           ClientVariant = "GetCwd"
	ClientGetIdleStatus    ClientVariant = "GetIdleStatus"
	ClientScroll           ClientVariant = "Scroll"
	ClientClosePty         ClientVariant = "ClosePty"
)

type ClientMsg struct {
	Variant ClientVariant
	Payload any
}

type ClientHelloPayload struct {
	Version string `json:"version"`
}

type ClientSpawnPayload struct {
	TabID string      `json:"tab_id"`
	Rows  uint16      `json:"rows"`
	Cols  uint16      `json:"cols"`
	Cwd   *string     `json:"cwd,omitempty"`
	Shell *string     `json:"shell,omitempty"`
	Env   [][2]string `json:"env,omitempty"`
}

type ClientInputPayload struct {
	TabID string    `json:"tab_id"`
	Data  ByteArray `json:"data"`
}

type ClientPastePayload struct {
	TabID string `json:"tab_id"`
	Data  string `json:"data"`
}

type ClientGetInputModesPayload struct {
	TabID string `json:"tab_id"`
}

type ClientResizePayload struct {
	TabID string `json:"tab_id"`
	Rows  uint16 `json:"rows"`
	Cols  uint16 `json:"cols"`
}

type ClientGetScreenPayload struct {
	TabID string `json:"tab_id"`
	Rows  uint16 `json:"rows"`
	Cols  uint16 `json:"cols"`
}

type ClientSubscribePayload struct {
	TabID string `json:"tab_id"`
	Rows  uint16 `json:"rows"`
	Cols  uint16 `json:"cols"`
}

type ClientScrollPayload struct {
	TabID string `json:"tab_id"`
	Delta int32  `json:"delta"`
}

type ClientClosePtyPayload struct {
	TabID string `json:"tab_id"`
}

type ClientGetCwdPayload struct {
	TabID string `json:"tab_id"`
}

func (m ClientMsg) MarshalJSON() ([]byte, error) {
	if m.Variant == "" {
		return nil, ErrInvalidEnvelope
	}
	if m.Payload == nil {
		m.Payload = map[string]any{}
	}
	env := map[string]any{string(m.Variant): m.Payload}
	return json.Marshal(env)
}

func (m *ClientMsg) UnmarshalJSON(data []byte) error {
	var env map[string]json.RawMessage
	if err := json.Unmarshal(data, &env); err != nil {
		return err
	}
	if len(env) != 1 {
		return ErrInvalidEnvelope
	}

	for k, v := range env {
		m.Variant = ClientVariant(k)
		switch m.Variant {
		case ClientHello:
			var p ClientHelloPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ClientSpawn:
			var p ClientSpawnPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ClientInput:
			var p ClientInputPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ClientPaste:
			var p ClientPastePayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ClientGetInputModes:
			var p ClientGetInputModesPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ClientResize:
			var p ClientResizePayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ClientGetScreen:
			var p ClientGetScreenPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ClientGetProcessStatus:
			m.Payload = struct{}{}
		case ClientSubscribe:
			var p ClientSubscribePayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ClientGetCwd:
			var p ClientGetCwdPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ClientGetIdleStatus:
			m.Payload = struct{}{}
		case ClientScroll:
			var p ClientScrollPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ClientClosePty:
			var p ClientClosePtyPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		default:
			return fmt.Errorf("%w: %s", ErrUnknownVariant, k)
		}
	}
	return nil
}

type ServerVariant string

const (
	ServerWelcome       ServerVariant = "Welcome"
	ServerError         ServerVariant = "Error"
	ServerInputModes    ServerVariant = "InputModes"
	ServerUpdateStatus  ServerVariant = "UpdateStatus"
	ServerProcessStatus ServerVariant = "ProcessStatus"
	ServerCwd           ServerVariant = "Cwd"
)

type ServerMsg struct {
	Variant ServerVariant
	Payload any
}

type ServerWelcomePayload struct {
	Version string `json:"version"`
}

type ServerErrorPayload struct {
	Message string `json:"message"`
}

type ServerInputModesPayload struct {
	Mouse          bool `json:"mouse"`
	BracketedPaste bool `json:"bracketed_paste"`
}

type ServerUpdateStatusPayload struct {
	Latest *string `json:"latest"`
}

type ServerProcessStatusPayload struct {
	Tabs [][2]any `json:"tabs"`
}

type ServerCwdPayload struct {
	TabID string  `json:"tab_id"`
	Path  *string `json:"path"`
}

func (m ServerMsg) MarshalJSON() ([]byte, error) {
	if m.Variant == "" {
		return nil, ErrInvalidEnvelope
	}
	if m.Payload == nil {
		m.Payload = map[string]any{}
	}
	env := map[string]any{string(m.Variant): m.Payload}
	return json.Marshal(env)
}

func (m *ServerMsg) UnmarshalJSON(data []byte) error {
	var env map[string]json.RawMessage
	if err := json.Unmarshal(data, &env); err != nil {
		return err
	}
	if len(env) != 1 {
		return ErrInvalidEnvelope
	}

	for k, v := range env {
		m.Variant = ServerVariant(k)
		switch m.Variant {
		case ServerWelcome:
			var p ServerWelcomePayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ServerError:
			var p ServerErrorPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ServerInputModes:
			var p ServerInputModesPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ServerUpdateStatus:
			var p ServerUpdateStatusPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ServerProcessStatus:
			var p ServerProcessStatusPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		case ServerCwd:
			var p ServerCwdPayload
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			m.Payload = p
		default:
			return fmt.Errorf("%w: %s", ErrUnknownVariant, k)
		}
	}
	return nil
}
