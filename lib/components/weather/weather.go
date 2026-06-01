package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gregdel/homed/lib/components"
	"github.com/gregdel/homed/lib/components/common"
	"github.com/gregdel/homed/lib/config"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultEndpoint = "https://wttr.in"
	defaultInterval = time.Hour
	requestTimeout  = 10 * time.Second
)

func init() {
	components.Register(components.TypeWeather, New)
}

// Params represents the weather component params.
type Params struct {
	Location *config.Location `yaml:"location"`
	Interval string           `yaml:"interval"`
}

// Data represents current weather data.
type Data struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
}

// Snapshot represents the HTTP and websocket JSON values for Weather.
type Snapshot struct {
	common.SnapshotBase
	Data
}

// Weather fetches current weather data from wttr.in.
type Weather struct {
	common.Component

	mu       sync.RWMutex
	Params   Params
	interval time.Duration
	client   *http.Client
	endpoint string
	Data
}

// New returns a new weather component.
func New() components.Component {
	return &Weather{
		client:   &http.Client{Timeout: requestTimeout},
		endpoint: defaultEndpoint,
		interval: defaultInterval,
	}
}

// Type implements the Component interface.
func (w *Weather) Type() components.Type {
	return components.TypeWeather
}

// Run implements the RunnableComponent interface.
func (w *Weather) Run(ctx context.Context, logger *slog.Logger, _ *components.Components) error {
	if err := w.setup(); err != nil {
		return err
	}

	log := w.LoggerWithFields(logger)
	w.fetchAndNotify(ctx, log)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("stopping")
			return nil
		case <-ticker.C:
			w.fetchAndNotify(ctx, log)
		}
	}
}

func (w *Weather) setup() error {
	if err := w.YAMLParams.Decode(&w.Params); err != nil {
		return err
	}

	return w.normalizeParams()
}

func (w *Weather) normalizeParams() error {
	if w.client == nil {
		w.client = &http.Client{Timeout: requestTimeout}
	}
	if w.endpoint == "" {
		w.endpoint = defaultEndpoint
	}

	if w.Params.Location == nil {
		return fmt.Errorf("components: weather: missing location")
	}
	if w.Params.Location.Latitude < -90 || w.Params.Location.Latitude > 90 {
		return fmt.Errorf("components: weather: latitude must be between -90 and 90")
	}
	if w.Params.Location.Longitude < -180 || w.Params.Location.Longitude > 180 {
		return fmt.Errorf("components: weather: longitude must be between -180 and 180")
	}

	if w.Params.Interval == "" {
		w.interval = defaultInterval
		return nil
	}

	interval, err := time.ParseDuration(w.Params.Interval)
	if err != nil {
		return fmt.Errorf("components: weather: invalid interval: %w", err)
	}
	if interval <= 0 {
		return fmt.Errorf("components: weather: interval must be positive")
	}

	w.interval = interval
	return nil
}

func (w *Weather) fetchAndNotify(ctx context.Context, log *slog.Logger) {
	if err := w.fetchAndUpdate(ctx); err != nil {
		log.Warn("failed to fetch weather", slog.Any("error", err))
		if device := w.Device(); device != nil && device.IsOnline() {
			device.SetOnline(false)
			w.Notify()
		}
		return
	}

	if device := w.Device(); device != nil {
		device.SetOnline(true)
	}

	if err := w.PostUpdate(); err != nil {
		log.Warn("failed to update weather component", slog.Any("error", err))
	}
}

func (w *Weather) fetchAndUpdate(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.requestURL(), nil)
	if err != nil {
		return err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("wttr.in returned status %s", resp.Status)
	}

	var response wttrResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	data, err := response.currentData()
	if err != nil {
		return err
	}

	w.setData(data)
	return nil
}

func (w *Weather) requestURL() string {
	u, err := url.Parse(w.endpoint)
	if err != nil {
		return w.endpoint
	}

	location := strconv.FormatFloat(w.Params.Location.Latitude, 'f', -1, 64) + "," +
		strconv.FormatFloat(w.Params.Location.Longitude, 'f', -1, 64)
	u.Path = strings.TrimRight(u.Path, "/") + "/" + location
	q := u.Query()
	q.Set("format", "j1")
	u.RawQuery = q.Encode()

	return u.String()
}

// Update implements the Component interface.
func (w *Weather) Update(value []byte) error {
	var data Data
	if err := json.Unmarshal(value, &data); err != nil {
		return err
	}

	w.setData(data)
	return nil
}

func (w *Weather) setData(data Data) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Data = data
}

func (w *Weather) dataSnapshot() Data {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.Data
}

// Snapshot returns the current HTTP and websocket JSON values shape.
func (w *Weather) ValuesSnapshot() any {
	return Snapshot{
		SnapshotBase: w.SnapshotBase(),
		Data:         w.dataSnapshot(),
	}
}

// Collectors implements the Component interface.
func (w *Weather) Collectors(labels prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{
		components.GaugeCollector("temperature", labels,
			func() float64 { return w.dataSnapshot().Temperature }),
		components.GaugeCollector("humidity", labels,
			func() float64 { return w.dataSnapshot().Humidity }),
		components.GaugeCollector("pressure", labels,
			func() float64 { return w.dataSnapshot().Pressure }),
	}
}

type wttrResponse struct {
	CurrentCondition []wttrCurrentCondition `json:"current_condition"`
}

type wttrCurrentCondition struct {
	Temperature string `json:"temp_C"`
	Humidity    string `json:"humidity"`
	Pressure    string `json:"pressure"`
}

func (r wttrResponse) currentData() (Data, error) {
	if len(r.CurrentCondition) == 0 {
		return Data{}, fmt.Errorf("wttr.in response missing current_condition")
	}

	current := r.CurrentCondition[0]
	data := Data{}

	var err error
	if data.Temperature, err = parseWttrFloat("temp_C", current.Temperature); err != nil {
		return Data{}, err
	}
	if data.Humidity, err = parseWttrFloat("humidity", current.Humidity); err != nil {
		return Data{}, err
	}
	if data.Pressure, err = parseWttrFloat("pressure", current.Pressure); err != nil {
		return Data{}, err
	}

	return data, nil
}

func parseWttrFloat(field, value string) (float64, error) {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("wttr.in field %s: %w", field, err)
	}

	return parsed, nil
}
