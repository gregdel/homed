package components

import (
	"log/slog"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gregdel/homed/lib/config"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestAvailabilityAwareCollectorsKeepUntrackedDeviceMetrics(t *testing.T) {
	device := NewDevice("sensor", "office")
	collector := GaugeCollector("test_untracked", prometheus.Labels{"id": "untracked"}, func() float64 {
		return 21.5
	})

	registry := prometheus.NewRegistry()
	for _, collector := range availabilityAwareCollectors(device, TypeGenericSensor, []prometheus.Collector{collector}) {
		if err := registry.Register(collector); err != nil {
			t.Fatalf("failed to register collector: %s", err)
		}
	}

	metric := gatherMetric(t, registry, "homed_test_untracked")
	if metric == nil {
		t.Fatal("expected metric to be exported")
	}
	if got := metric.GetGauge().GetValue(); got != 21.5 {
		t.Fatalf("unexpected metric value: got %f, want %f", got, 21.5)
	}
}

func TestAvailabilityAwareCollectorsSuppressOfflineTrackedTelemetry(t *testing.T) {
	device := NewDevice("sensor", "office")
	device.SetAvailabilityTracked()

	called := false
	collector := GaugeCollector("test_offline", prometheus.Labels{"id": "offline"}, func() float64 {
		called = true
		return 21.5
	})

	registry := prometheus.NewRegistry()
	for _, collector := range availabilityAwareCollectors(device, TypeGenericSensor, []prometheus.Collector{collector}) {
		if err := registry.Register(collector); err != nil {
			t.Fatalf("failed to register collector: %s", err)
		}
	}

	if metric := gatherMetric(t, registry, "homed_test_offline"); metric != nil {
		t.Fatal("expected offline telemetry metric to be suppressed")
	}
	if called {
		t.Fatal("expected telemetry callback not to be called while offline")
	}

	device.SetOnline(true)
	metric := gatherMetric(t, registry, "homed_test_offline")
	if metric == nil {
		t.Fatal("expected metric to be exported once the device is online")
	}
	if got := metric.GetGauge().GetValue(); got != 21.5 {
		t.Fatalf("unexpected metric value: got %f, want %f", got, 21.5)
	}
}

func TestAvailabilityAwareCollectorsKeepDeviceStatusWhenOffline(t *testing.T) {
	device := NewDevice("sensor", "office")
	device.SetAvailabilityTracked()
	collector := GaugeCollector("test_device_status", prometheus.Labels{"id": "status"}, func() float64 {
		return 0
	})

	registry := prometheus.NewRegistry()
	for _, collector := range availabilityAwareCollectors(device, TypeDeviceStatus, []prometheus.Collector{collector}) {
		if err := registry.Register(collector); err != nil {
			t.Fatalf("failed to register collector: %s", err)
		}
	}

	metric := gatherMetric(t, registry, "homed_test_device_status")
	if metric == nil {
		t.Fatal("expected device status metric to be exported while offline")
	}
	if got := metric.GetGauge().GetValue(); got != 0 {
		t.Fatalf("unexpected metric value: got %f, want 0", got)
	}
}

func TestComponentsAddSuppressesExistingTelemetryWhenAvailabilityProviderIsAdded(t *testing.T) {
	telemetryType := Type("test_add_telemetry")
	setRegisteredComponent(t, telemetryType, func() Component {
		return &testComponent{
			componentType: telemetryType,
			collector: GaugeCollector("test_add_telemetry", prometheus.Labels{}, func() float64 {
				return 21.5
			}),
		}
	})
	setRegisteredComponent(t, TypeDeviceStatus, func() Component {
		return &testAvailabilityComponent{
			testComponent: testComponent{
				componentType: TypeDeviceStatus,
				collector: GaugeCollector("test_add_device_status", prometheus.Labels{}, func() float64 {
					return 0
				}),
			},
		}
	})

	registry := prometheus.NewRegistry()
	previousRegisterer := prometheus.DefaultRegisterer
	prometheus.DefaultRegisterer = registry
	t.Cleanup(func() {
		prometheus.DefaultRegisterer = previousRegisterer
	})

	components := New("")
	if _, err := components.Add(config.Component{Type: string(telemetryType)}, "office", "sensor"); err != nil {
		t.Fatalf("failed to add telemetry component: %s", err)
	}
	if metric := gatherMetric(t, registry, "homed_test_add_telemetry"); metric == nil {
		t.Fatal("expected telemetry metric before availability provider is added")
	}

	if _, err := components.Add(config.Component{Type: string(TypeDeviceStatus)}, "office", "sensor"); err != nil {
		t.Fatalf("failed to add availability component: %s", err)
	}
	if metric := gatherMetric(t, registry, "homed_test_add_telemetry"); metric != nil {
		t.Fatal("expected existing telemetry metric to be suppressed while tracked device is offline")
	}
	if metric := gatherMetric(t, registry, "homed_test_add_device_status"); metric == nil {
		t.Fatal("expected device status metric to remain visible while tracked device is offline")
	}
}

func gatherMetric(t *testing.T, registry *prometheus.Registry, name string) *dto.Metric {
	t.Helper()

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %s", err)
	}

	for _, family := range families {
		if family.GetName() != name {
			continue
		}
		if len(family.Metric) == 0 {
			return nil
		}
		return family.Metric[0]
	}

	return nil
}

func setRegisteredComponent(t *testing.T, componentType Type, fn func() Component) {
	t.Helper()

	name := string(componentType)
	previous, hadPrevious := registeredComponents[name]
	registeredComponents[name] = fn
	t.Cleanup(func() {
		if hadPrevious {
			registeredComponents[name] = previous
			return
		}
		delete(registeredComponents, name)
	})
}

type testComponent struct {
	componentType Type
	id            string
	device        *Device
	config        *config.Component
	collector     prometheus.Collector
	mqttClient    mqtt.Client
}

func (c *testComponent) Type() Type { return c.componentType }

func (c *testComponent) FriendlyName() string { return "" }

func (c *testComponent) Device() *Device { return c.device }

func (c *testComponent) SetDevice(device *Device) { c.device = device }

func (c *testComponent) Config() *config.Component { return c.config }

func (c *testComponent) SetConfig(config *config.Component) { c.config = config }

func (c *testComponent) LoggerWithFields(logger *slog.Logger) *slog.Logger { return logger }

func (c *testComponent) Update([]byte) error { return nil }

func (c *testComponent) PostUpdate() error { return nil }

func (c *testComponent) Collectors(prometheus.Labels) []prometheus.Collector {
	return []prometheus.Collector{c.collector}
}

func (c *testComponent) ReadOnly() bool { return true }

func (c *testComponent) Internal() bool { return false }

func (c *testComponent) WriteCommand([]byte) error { return ErrNotImplemented }

func (c *testComponent) ExecCommand([]byte) error { return ErrNotImplemented }

func (c *testComponent) PublishToStateTopic([]byte) error { return ErrNotImplemented }

func (c *testComponent) SetMQTTClient(client mqtt.Client) { c.mqttClient = client }

func (c *testComponent) MQTTClient() mqtt.Client { return c.mqttClient }

func (c *testComponent) ID() string { return c.id }

func (c *testComponent) SetID(id string) { c.id = id }

func (c *testComponent) Subscribe(string, chan Event) {}

func (c *testComponent) Notify() {}

type testAvailabilityComponent struct {
	testComponent
}

func (c *testAvailabilityComponent) ProvidesAvailability() {}
