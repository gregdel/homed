package httpd

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"github.com/gregdel/homed/lib/components"
	_ "github.com/gregdel/homed/lib/components/script"
	"github.com/gregdel/homed/lib/config"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"gopkg.in/yaml.v3"
)

func TestUpdateComponentPublishesInternalScriptCommandToMQTT(t *testing.T) {
	var params yaml.Node
	if err := yaml.Unmarshal([]byte("command: /bin/sh\nargs: [\"-c\", \"exit 0\"]\ntimeout: 2s"), &params); err != nil {
		t.Fatalf("failed to decode params: %s", err)
	}

	inventory := components.New("")
	component, err := inventory.Add(config.Component{
		ID:           "scan_page",
		Type:         string(components.TypeScript),
		FriendlyName: "Scan page",
		Internal:     true,
		StateTopic:   "home/homed/scripts/scan_page/state",
		CommandTopic: "home/homed/scripts/scan_page/command",
		Params:       params,
	}, "none", "automations")
	if err != nil {
		t.Fatalf("failed to add script: %s", err)
	}
	client := &fakeMQTTClient{connected: true}
	component.SetMQTTClient(client)

	h := &httpd{
		components: inventory,
		render:     render.New(),
	}

	req := httptest.NewRequest(http.MethodPut, "/components/scan_page", strings.NewReader("{}"))
	rec := httptest.NewRecorder()

	h.updateComponent(rec, req, httprouter.Params{
		{Key: "id", Value: "scan_page"},
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"status":"success"`) {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}
	if client.topic != "home/homed/scripts/scan_page/command" {
		t.Fatalf("unexpected topic: got %q", client.topic)
	}
	if string(client.payload) != "{}" {
		t.Fatalf("unexpected payload: got %q", string(client.payload))
	}
}

func TestWebsocketEventsRegistersAndUnregistersClient(t *testing.T) {
	h := &httpd{
		logger:     slogDiscard(),
		websockets: map[*websocket.Conn]*websocketClient{},
		render:     render.New(),
	}

	router := httprouter.New()
	router.GET("/events", h.websocketEvents)
	server := httptest.NewServer(router)
	defer server.Close()

	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse test server URL: %s", err)
	}
	u.Scheme = "ws"
	u.Path = "/events"

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("failed to connect websocket: %s", err)
	}

	eventuallyWebsocketCount(t, h, 1)

	if err := conn.Close(); err != nil {
		t.Fatalf("failed to close websocket: %s", err)
	}

	eventuallyWebsocketCount(t, h, 0)
}

func eventuallyWebsocketCount(t *testing.T, h *httpd, want int) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		h.mu.RLock()
		got := len(h.websockets)
		h.mu.RUnlock()

		if got == want {
			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	h.mu.RLock()
	got := len(h.websockets)
	h.mu.RUnlock()
	t.Fatalf("unexpected websocket count: got %d, want %d", got, want)
}

func slogDiscard() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type fakeMQTTClient struct {
	connected bool
	topic     string
	payload   []byte
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

func (c *fakeMQTTClient) Publish(topic string, _ byte, _ bool, payload any) mqtt.Token {
	c.topic = topic
	switch p := payload.(type) {
	case []byte:
		c.payload = append([]byte(nil), p...)
	case string:
		c.payload = []byte(p)
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
