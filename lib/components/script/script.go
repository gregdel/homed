package script

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultTimeout = 5 * time.Minute

	StatusIdle    = "idle"
	StatusRunning = "running"
	StatusSuccess = "success"
	StatusError   = "error"
)

func init() {
	components.Register(components.TypeScript, New)
}

type Params struct {
	Command string        `yaml:"command"`
	Args    []string      `yaml:"args"`
	Timeout time.Duration `yaml:"timeout"`
}

type Script struct {
	common.Component

	mu sync.RWMutex

	params Params

	running        bool
	status         string
	lastStartedAt  *time.Time
	lastFinishedAt *time.Time
	lastDuration   time.Duration
	lastExitCode   *int
	lastError      string
}

type Snapshot struct {
	common.SnapshotBase
	State
}

type State struct {
	Status         string     `json:"status"`
	LastStartedAt  *time.Time `json:"last_started_at"`
	LastFinishedAt *time.Time `json:"last_finished_at"`
	LastDurationMS int64      `json:"last_duration_ms"`
	LastExitCode   *int       `json:"last_exit_code"`
	LastError      string     `json:"last_error"`
}

func New() components.Component {
	s := &Script{}
	s.status = StatusIdle
	return s
}

func (s *Script) Type() components.Type {
	return components.TypeScript
}

func (s *Script) ValidateConfig() error {
	if !s.Internal() {
		return errors.New("script: internal must be true")
	}
	if s.StateTopic == "" {
		return errors.New("script: state_topic is required")
	}
	if s.CommandTopic == "" {
		return errors.New("script: command_topic is required")
	}

	params := Params{Timeout: defaultTimeout}
	if err := s.YAMLParams.Decode(&params); err != nil {
		return fmt.Errorf("script: failed to decode params: %w", err)
	}
	if params.Command == "" {
		return errors.New("script: params.command is required")
	}
	if params.Timeout <= 0 {
		return errors.New("script: params.timeout must be greater than zero")
	}

	s.mu.Lock()
	s.params = params
	s.mu.Unlock()

	device := s.Device()
	if device != nil {
		device.SetOnline(true)
	}

	return nil
}

func (s *Script) Update(data []byte) error {
	state := State{}
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	if state.Status == "" {
		state.Status = StatusIdle
	}

	s.mu.Lock()
	if state.Status == StatusRunning && (!s.running || s.status != StatusRunning) {
		state.Status = StatusIdle
	}
	s.status = state.Status
	s.lastStartedAt = state.LastStartedAt
	s.lastFinishedAt = state.LastFinishedAt
	s.lastDuration = time.Duration(state.LastDurationMS) * time.Millisecond
	s.lastExitCode = state.LastExitCode
	s.lastError = state.LastError
	s.mu.Unlock()

	return nil
}

func (s *Script) ExecCommand([]byte) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return components.ErrOperatingInProgress
	}

	params := s.params
	startedAt := time.Now()
	s.running = true
	s.status = StatusRunning
	s.lastStartedAt = &startedAt
	s.lastFinishedAt = nil
	s.lastDuration = 0
	s.lastExitCode = nil
	s.lastError = ""
	state := s.stateSnapshotLocked()
	s.mu.Unlock()

	if err := s.publishStateSnapshot(state); err != nil {
		s.mu.Lock()
		s.running = false
		s.status = StatusIdle
		s.lastStartedAt = nil
		s.mu.Unlock()
		return err
	}

	go s.run(params, startedAt)
	return nil
}

func (s *Script) run(params Params, startedAt time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), params.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, params.Command, params.Args...)
	output, err := cmd.CombinedOutput()

	finishedAt := time.Now()
	exitCode := commandExitCode(err)
	status := StatusSuccess
	lastError := ""
	if err != nil {
		status = StatusError
		if ctx.Err() == context.DeadlineExceeded {
			lastError = fmt.Sprintf("timeout after %s", params.Timeout)
		} else {
			lastError = err.Error()
		}
	}

	s.mu.Lock()
	s.status = status
	s.lastFinishedAt = &finishedAt
	s.lastDuration = finishedAt.Sub(startedAt)
	s.lastExitCode = exitCode
	s.lastError = lastError
	state := s.stateSnapshotLocked()
	s.mu.Unlock()

	logger := s.LoggerWithFields(slog.Default())
	if err != nil {
		logger.Warn("script command failed",
			slog.String("command", params.Command),
			slog.Any("args", params.Args),
			slog.Any("error", err),
			slog.String("output", string(output)),
		)
	} else {
		logger.Info("script command finished",
			slog.String("command", params.Command),
			slog.Any("args", params.Args),
			slog.String("output", string(output)),
		)
	}

	if err := s.publishStateSnapshot(state); err != nil {
		logger.Warn("failed to publish script update", slog.Any("error", err))
	}

	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
}

func commandExitCode(err error) *int {
	if err == nil {
		code := 0
		return &code
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return nil
	}

	if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
		code := status.ExitStatus()
		return &code
	}

	code := exitErr.ExitCode()
	return &code
}

func (s *Script) ValuesSnapshot() any {
	state := s.stateSnapshot()

	return Snapshot{
		SnapshotBase: s.SnapshotBase(),
		State:        state,
	}
}

func (s *Script) stateSnapshot() State {
	s.mu.RLock()
	state := s.stateSnapshotLocked()
	s.mu.RUnlock()

	return state
}

func (s *Script) stateSnapshotLocked() State {
	status := s.status
	if status == "" {
		status = StatusIdle
	}

	return State{
		Status:         status,
		LastStartedAt:  copyTimePtr(s.lastStartedAt),
		LastFinishedAt: copyTimePtr(s.lastFinishedAt),
		LastDurationMS: s.lastDuration.Milliseconds(),
		LastExitCode:   copyIntPtr(s.lastExitCode),
		LastError:      s.lastError,
	}
}

func (s *Script) publishState() error {
	return s.publishStateSnapshot(s.stateSnapshot())
}

func (s *Script) publishStateSnapshot(state State) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return s.PublishToStateTopic(data)
}

func copyTimePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	copied := *t
	return &copied
}

func copyIntPtr(i *int) *int {
	if i == nil {
		return nil
	}
	copied := *i
	return &copied
}

func (s *Script) Collectors(prometheus.Labels) []prometheus.Collector {
	return nil
}
