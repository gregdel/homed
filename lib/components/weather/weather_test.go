package weather

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/config"
	"gopkg.in/yaml.v3"
)

func TestFetchAndUpdateRequestsWttrLocationAndParsesCurrentConditions(t *testing.T) {
	var requestPath, requestQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath = r.URL.Path
		requestQuery = r.URL.RawQuery

		_, _ = w.Write([]byte(`{
			"current_condition": [{
				"temp_C": "12",
				"FeelsLikeC": "10",
				"humidity": "88",
				"pressure": "1004",
				"precipMM": "0.2",
				"cloudcover": "75",
				"visibility": "9",
				"observation_time": "09:30 AM",
				"weatherDesc": [{"value": "Light rain"}],
				"windspeedKmph": "11"
			}]
		}`))
	}))
	defer server.Close()

	w := newTestWeather(server.URL)
	if err := w.fetchAndUpdate(context.Background()); err != nil {
		t.Fatalf("fetchAndUpdate returned error: %s", err)
	}

	if requestPath != "/50.62955091282183,3.055823770700181" {
		t.Fatalf("unexpected request path: got %q", requestPath)
	}
	if requestQuery != "format=j1" {
		t.Fatalf("unexpected request query: got %q", requestQuery)
	}

	got := w.dataSnapshot()
	want := Data{
		Temperature: 12,
		Humidity:    88,
		Pressure:    1004,
	}
	if got != want {
		t.Fatalf("unexpected weather data: got %+v, want %+v", got, want)
	}
}

func TestSnapshotOnlyExposesClimateFields(t *testing.T) {
	w := newTestWeather("http://example.test")
	w.setData(Data{
		Temperature: 12,
		Humidity:    88,
		Pressure:    1004,
	})

	data, err := json.Marshal(components.NewComponentJSON(w, true))
	if err != nil {
		t.Fatalf("failed to marshal weather snapshot: %s", err)
	}

	for _, field := range []string{
		"wind",
		"feels_like",
		"precipitation",
		"cloud_cover",
		"visibility",
		"description",
		"observation_time",
	} {
		if strings.Contains(string(data), field) {
			t.Fatalf("snapshot should not contain %q: %s", field, data)
		}
	}
	if !strings.Contains(string(data), `"temperature":12`) {
		t.Fatalf("snapshot missing temperature data: %s", data)
	}
}

func TestFetchAndUpdateErrorsDoNotCorruptExistingState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	component := newTestWeather(server.URL)
	existing := Data{Temperature: 18, Humidity: 40, Pressure: 1012}
	component.setData(existing)

	if err := component.fetchAndUpdate(context.Background()); err == nil {
		t.Fatal("expected fetchAndUpdate to return an error")
	}

	if got := component.dataSnapshot(); got != existing {
		t.Fatalf("state changed after failed fetch: got %+v, want %+v", got, existing)
	}
}

func TestFetchAndNotifyMarksOnlineDeviceOfflineOnError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	component := newTestWeather(server.URL)
	component.Device().SetOnline(true)
	component.setData(Data{Temperature: 18, Humidity: 40, Pressure: 1012})

	updatedAt := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	component.UpdatedAt.Store(&updatedAt)

	events := components.NewEventChannel()
	component.Subscribe("test", events)

	component.fetchAndNotify(context.Background(), slog.New(slog.NewTextHandler(io.Discard, nil)))

	if component.Device().IsOnline() {
		t.Fatal("expected device to be marked offline after failed fetch")
	}
	if got := component.dataSnapshot(); got != (Data{Temperature: 18, Humidity: 40, Pressure: 1012}) {
		t.Fatalf("state changed after failed fetch: got %+v", got)
	}
	if got := component.UpdatedAt.Load(); got == nil || !got.Equal(updatedAt) {
		t.Fatalf("updated_at changed after failed fetch: got %v, want %v", got, updatedAt)
	}

	select {
	case event := <-events:
		if event.ID != component.ID() {
			t.Fatalf("unexpected event id: got %q, want %q", event.ID, component.ID())
		}
	default:
		t.Fatal("expected failed fetch to notify subscribers")
	}
}

func TestFetchAndNotifyDoesNotNotifyWhenAlreadyOffline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	component := newTestWeather(server.URL)
	events := components.NewEventChannel()
	component.Subscribe("test", events)

	component.fetchAndNotify(context.Background(), slog.New(slog.NewTextHandler(io.Discard, nil)))

	select {
	case event := <-events:
		t.Fatalf("unexpected event while already offline: %+v", event)
	default:
	}
}

func TestFetchAndUpdateReturnsErrorOnInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"current_condition":[{"temp_C":"oops"}]}`))
	}))
	defer server.Close()

	component := newTestWeather(server.URL)
	if err := component.fetchAndUpdate(context.Background()); err == nil {
		t.Fatal("expected fetchAndUpdate to return an error")
	}
}

func TestNormalizeParamsDefaultsInterval(t *testing.T) {
	w := &Weather{Params: Params{Location: testLocation()}}
	if err := w.normalizeParams(); err != nil {
		t.Fatalf("normalizeParams returned error: %s", err)
	}

	if w.interval != defaultInterval {
		t.Fatalf("unexpected interval: got %s, want %s", w.interval, defaultInterval)
	}
}

func TestNormalizeParamsRejectsInvalidIntervals(t *testing.T) {
	tests := []string{"bad", "0s", "-1h"}

	for _, interval := range tests {
		t.Run(interval, func(t *testing.T) {
			w := &Weather{Params: Params{Location: testLocation(), Interval: interval}}
			if err := w.normalizeParams(); err == nil {
				t.Fatal("expected interval validation error")
			}
		})
	}
}

func TestNormalizeParamsRejectsMissingLocation(t *testing.T) {
	w := &Weather{}
	if err := w.normalizeParams(); err == nil {
		t.Fatal("expected missing location validation error")
	}
}

func TestNormalizeParamsRejectsOutOfRangeLocation(t *testing.T) {
	tests := []struct {
		name     string
		location config.Location
	}{
		{
			name:     "latitude below range",
			location: config.Location{Latitude: -91, Longitude: 0},
		},
		{
			name:     "latitude above range",
			location: config.Location{Latitude: 91, Longitude: 0},
		},
		{
			name:     "longitude below range",
			location: config.Location{Latitude: 0, Longitude: -181},
		},
		{
			name:     "longitude above range",
			location: config.Location{Latitude: 0, Longitude: 181},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := &Weather{Params: Params{Location: &tc.location}}
			if err := w.normalizeParams(); err == nil {
				t.Fatal("expected location validation error")
			}
		})
	}
}

func TestNormalizeParamsAcceptsExplicitZeroLocation(t *testing.T) {
	location := config.Location{Latitude: 0, Longitude: 0}
	w := &Weather{Params: Params{Location: &location}}
	if err := w.normalizeParams(); err != nil {
		t.Fatalf("normalizeParams returned error: %s", err)
	}
}

func TestSetupDecodesParams(t *testing.T) {
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(`location:
  latitude: 1.25
  longitude: 2.5
interval: 30m
`), &node); err != nil {
		t.Fatalf("failed to build yaml node: %s", err)
	}

	w := New().(*Weather)
	w.YAMLParams = node
	if err := w.setup(); err != nil {
		t.Fatalf("setup returned error: %s", err)
	}

	if w.Params.Location == nil {
		t.Fatal("expected location to be decoded")
	}
	if w.Params.Location.Latitude != 1.25 || w.Params.Location.Longitude != 2.5 {
		t.Fatalf("unexpected location: %+v", w.Params.Location)
	}
	if w.interval != 30*time.Minute {
		t.Fatalf("unexpected interval: got %s", w.interval)
	}
}

func newTestWeather(endpoint string) *Weather {
	w := New().(*Weather)
	w.endpoint = endpoint
	w.Params.Location = testLocation()
	w.SetID("weather")
	w.SetDevice(components.NewDevice("weather", "outside"))
	return w
}

func testLocation() *config.Location {
	return &config.Location{
		Latitude:  50.62955091282183,
		Longitude: 3.055823770700181,
	}
}
