package script

import (
	"encoding/json"
	"runtime"
	"sync"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
	"gopkg.in/yaml.v3"
)

func TestValidateConfigRequiresInternal(t *testing.T) {
	s := newConfiguredScript(t, "command: /bin/true")
	cfg := s.Config()
	cfg.Internal = false
	s.SetConfig(cfg)

	if err := s.ValidateConfig(); err == nil {
		t.Fatal("expected missing internal error")
	}
}

func TestValidateConfigRequiresStateTopic(t *testing.T) {
	s := newConfiguredScript(t, "command: /bin/true")
	cfg := s.Config()
	cfg.StateTopic = ""
	s.SetConfig(cfg)

	if err := s.ValidateConfig(); err == nil {
		t.Fatal("expected missing state topic error")
	}
}

func TestValidateConfigRequiresCommandTopic(t *testing.T) {
	s := newConfiguredScript(t, "command: /bin/true")
	cfg := s.Config()
	cfg.CommandTopic = ""
	s.SetConfig(cfg)

	if err := s.ValidateConfig(); err == nil {
		t.Fatal("expected missing command topic error")
	}
}

func TestValidateConfigRequiresCommand(t *testing.T) {
	s := newConfiguredScript(t, "timeout: 1s")

	if err := s.ValidateConfig(); err == nil {
		t.Fatal("expected missing command error")
	}
}

func TestValidateConfigRejectsInvalidTimeout(t *testing.T) {
	s := newConfiguredScript(t, "command: /bin/true\ntimeout: 0s")

	if err := s.ValidateConfig(); err == nil {
		t.Fatal("expected invalid timeout error")
	}
}

func TestScriptPublishesRunningAndFinalState(t *testing.T) {
	s := validatedScript(t, "command: /bin/sh\nargs: [\"-c\", \"exit 0\"]\ntimeout: 2s")
	client := s.MQTTClient().(*fakeMQTTClient)

	if err := s.ExecCommand(nil); err != nil {
		t.Fatalf("failed to start script: %s", err)
	}

	waitForStatus(t, s, StatusSuccess)

	states := client.states()
	if len(states) < 2 {
		t.Fatalf("expected at least 2 state publishes, got %d", len(states))
	}
	if states[0].Status != StatusRunning {
		t.Fatalf("unexpected first state: got %q, want %q", states[0].Status, StatusRunning)
	}
	if states[len(states)-1].Status != StatusSuccess {
		t.Fatalf("unexpected final state: got %q, want %q", states[len(states)-1].Status, StatusSuccess)
	}
}

func TestScriptRunsCommandAsync(t *testing.T) {
	s := validatedScript(t, "command: /bin/sh\nargs: [\"-c\", \"exit 0\"]\ntimeout: 2s")

	if err := s.ExecCommand(nil); err != nil {
		t.Fatalf("failed to start script: %s", err)
	}

	waitForStatus(t, s, StatusSuccess)

	snapshot := s.ValuesSnapshot().(Snapshot)
	if snapshot.LastExitCode == nil || *snapshot.LastExitCode != 0 {
		t.Fatalf("unexpected exit code: %v", snapshot.LastExitCode)
	}
	if snapshot.LastDurationMS < 0 {
		t.Fatalf("unexpected negative duration: %d", snapshot.LastDurationMS)
	}
}

func TestScriptRejectsDuplicateRuns(t *testing.T) {
	s := validatedScript(t, "command: /bin/sh\nargs: [\"-c\", \"sleep 0.2\"]\ntimeout: 2s")

	if err := s.ExecCommand(nil); err != nil {
		t.Fatalf("failed to start script: %s", err)
	}
	if err := s.ExecCommand(nil); err != components.ErrOperatingInProgress {
		t.Fatalf("unexpected duplicate error: %v", err)
	}

	waitForStatus(t, s, StatusSuccess)
}

func TestScriptRecordsNonZeroExit(t *testing.T) {
	s := validatedScript(t, "command: /bin/sh\nargs: [\"-c\", \"exit 7\"]\ntimeout: 2s")

	if err := s.ExecCommand(nil); err != nil {
		t.Fatalf("failed to start script: %s", err)
	}

	waitForStatus(t, s, StatusError)

	snapshot := s.ValuesSnapshot().(Snapshot)
	if snapshot.LastExitCode == nil || *snapshot.LastExitCode != 7 {
		t.Fatalf("unexpected exit code: %v", snapshot.LastExitCode)
	}
	if snapshot.LastError == "" {
		t.Fatal("expected last error")
	}
}

func TestScriptRecordsTimeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell command uses unix sleep")
	}

	s := validatedScript(t, "command: /bin/sh\nargs: [\"-c\", \"sleep 1\"]\ntimeout: 10ms")

	if err := s.ExecCommand(nil); err != nil {
		t.Fatalf("failed to start script: %s", err)
	}

	waitForStatus(t, s, StatusError)

	snapshot := s.ValuesSnapshot().(Snapshot)
	if snapshot.LastError == "" {
		t.Fatal("expected timeout error")
	}
}

func TestScriptUpdateDecodesState(t *testing.T) {
	s := validatedScript(t, "command: /bin/sh\nargs: [\"-c\", \"exit 0\"]\ntimeout: 2s")
	exitCode := 12
	startedAt := time.Now().Add(-time.Second)
	finishedAt := time.Now()
	state := State{
		Status:         StatusError,
		LastStartedAt:  &startedAt,
		LastFinishedAt: &finishedAt,
		LastDurationMS: 123,
		LastExitCode:   &exitCode,
		LastError:      "failed",
	}
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("failed to marshal state: %s", err)
	}

	if err := s.Update(data); err != nil {
		t.Fatalf("failed to update script: %s", err)
	}

	snapshot := s.ValuesSnapshot().(Snapshot)
	if snapshot.Status != StatusError {
		t.Fatalf("unexpected status: got %q, want %q", snapshot.Status, StatusError)
	}
	if snapshot.LastDurationMS != 123 {
		t.Fatalf("unexpected duration: got %d, want 123", snapshot.LastDurationMS)
	}
	if snapshot.LastExitCode == nil || *snapshot.LastExitCode != exitCode {
		t.Fatalf("unexpected exit code: %v", snapshot.LastExitCode)
	}
	if snapshot.LastError != "failed" {
		t.Fatalf("unexpected error: %q", snapshot.LastError)
	}
}

func TestScriptUpdateNormalizesStaleRunningState(t *testing.T) {
	s := validatedScript(t, "command: /bin/sh\nargs: [\"-c\", \"exit 0\"]\ntimeout: 2s")
	data, err := json.Marshal(State{Status: StatusRunning})
	if err != nil {
		t.Fatalf("failed to marshal state: %s", err)
	}

	if err := s.Update(data); err != nil {
		t.Fatalf("failed to update script: %s", err)
	}

	snapshot := s.ValuesSnapshot().(Snapshot)
	if snapshot.Status != StatusIdle {
		t.Fatalf("unexpected status: got %q, want %q", snapshot.Status, StatusIdle)
	}
}

func TestScriptSnapshotMarshalsStatus(t *testing.T) {
	s := validatedScript(t, "command: /bin/sh\nargs: [\"-c\", \"exit 0\"]\ntimeout: 2s")

	data, err := json.Marshal(components.NewComponentJSON(s, false))
	if err != nil {
		t.Fatalf("failed to marshal component JSON: %s", err)
	}

	var got struct {
		Type   string `json:"type"`
		Values struct {
			ID             string `json:"id"`
			Status         string `json:"status"`
			LastDurationMS int64  `json:"last_duration_ms"`
		} `json:"values"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode component JSON: %s", err)
	}

	if got.Type != string(components.TypeScript) {
		t.Fatalf("unexpected type: got %q, want %q", got.Type, components.TypeScript)
	}
	if got.Values.ID != "test_script" {
		t.Fatalf("unexpected id: got %q", got.Values.ID)
	}
	if got.Values.Status != StatusIdle {
		t.Fatalf("unexpected status: got %q, want %q", got.Values.Status, StatusIdle)
	}
}

func newConfiguredScript(t *testing.T, paramsYAML string) *Script {
	t.Helper()

	var params yaml.Node
	if err := yaml.Unmarshal([]byte(paramsYAML), &params); err != nil {
		t.Fatalf("failed to decode yaml: %s", err)
	}

	s := New().(*Script)
	s.SetID("test_script")
	s.SetDevice(components.NewDevice("automations", "none"))
	s.SetConfig(&config.Component{
		FriendlyName: "Test Script",
		Internal:     true,
		StateTopic:   "home/homed/scripts/test_script/state",
		CommandTopic: "home/homed/scripts/test_script/command",
		Params:       params,
	})
	s.SetMQTTClient(&fakeMQTTClient{connected: true})
	return s
}

func validatedScript(t *testing.T, paramsYAML string) *Script {
	t.Helper()

	s := newConfiguredScript(t, paramsYAML)
	if err := s.ValidateConfig(); err != nil {
		t.Fatalf("failed to validate script config: %s", err)
	}
	return s
}

func waitForStatus(t *testing.T, s *Script, status string) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		snapshot := s.ValuesSnapshot().(Snapshot)
		if snapshot.Status == status {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("timed out waiting for status %q, got %q", status, s.ValuesSnapshot().(Snapshot).Status)
}

type fakeMQTTClient struct {
	mu        sync.Mutex
	connected bool
	payloads  [][]byte
}

func (c *fakeMQTTClient) IsConnected() bool {
	return c.connected
}

func (c *fakeMQTTClient) IsConnectionOpen() bool {
	return c.connected
}

func (c *fakeMQTTClient) Connect() mqtt.Token {
	return fakeToken{}
}

func (c *fakeMQTTClient) Disconnect(uint) {}

func (c *fakeMQTTClient) Publish(_ string, _ byte, _ bool, payload any) mqtt.Token {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch p := payload.(type) {
	case []byte:
		c.payloads = append(c.payloads, append([]byte(nil), p...))
	case string:
		c.payloads = append(c.payloads, []byte(p))
	}
	return fakeToken{}
}

func (c *fakeMQTTClient) Subscribe(string, byte, mqtt.MessageHandler) mqtt.Token {
	return fakeToken{}
}

func (c *fakeMQTTClient) SubscribeMultiple(map[string]byte, mqtt.MessageHandler) mqtt.Token {
	return fakeToken{}
}

func (c *fakeMQTTClient) Unsubscribe(...string) mqtt.Token {
	return fakeToken{}
}

func (c *fakeMQTTClient) AddRoute(string, mqtt.MessageHandler) {}

func (c *fakeMQTTClient) OptionsReader() mqtt.ClientOptionsReader {
	return mqtt.ClientOptionsReader{}
}

func (c *fakeMQTTClient) states() []State {
	c.mu.Lock()
	payloads := make([][]byte, 0, len(c.payloads))
	for _, payload := range c.payloads {
		payloads = append(payloads, append([]byte(nil), payload...))
	}
	c.mu.Unlock()

	states := []State{}
	for _, payload := range payloads {
		state := State{}
		if err := json.Unmarshal(payload, &state); err == nil {
			states = append(states, state)
		}
	}
	return states
}

type fakeToken struct{}

func (fakeToken) Wait() bool {
	return true
}

func (fakeToken) WaitTimeout(time.Duration) bool {
	return true
}

func (fakeToken) Done() <-chan struct{} {
	done := make(chan struct{})
	close(done)
	return done
}

func (fakeToken) Error() error {
	return nil
}
