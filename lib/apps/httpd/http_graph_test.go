package httpd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestQueryPrometheusGraph(t *testing.T) {
	tests := []struct {
		name         string
		rangeName    string
		expectedStep string
	}{
		{name: "one hour", rangeName: "1h", expectedStep: "60"},
		{name: "one day", rangeName: "24h", expectedStep: "300"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/query_range" {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				if got := r.URL.Query().Get("query"); !strings.Contains(got, `id="living_room"`) {
					t.Fatalf("query does not filter by component id: %s", got)
				}
				if got := r.URL.Query().Get("step"); got != tt.expectedStep {
					t.Fatalf("unexpected step: %s", got)
				}

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					"status": "success",
					"data": {
						"result": [{
							"metric": {
								"__name__": "homed_temperature_control_current",
								"component_type": "homed_temperature",
								"id": "living_room"
							},
							"values": [[1710000000, "19.5"], [1710000300, "20.1"]]
						}]
					}
				}`))
			}))
			defer server.Close()

			h := &httpd{
				httpClient:    server.Client(),
				prometheusURL: server.URL,
			}

			graph, err := h.queryPrometheusGraph("living_room", tt.rangeName, graphRanges[tt.rangeName])
			if err != nil {
				t.Fatalf("queryPrometheusGraph returned error: %s", err)
			}

			if graph.ComponentID != "living_room" {
				t.Fatalf("unexpected component id: %s", graph.ComponentID)
			}
			if graph.Range != tt.rangeName {
				t.Fatalf("unexpected range: %s", graph.Range)
			}
			if graph.StepSeconds != int64(graphRanges[tt.rangeName].step.Seconds()) {
				t.Fatalf("unexpected step seconds: %d", graph.StepSeconds)
			}
			if len(graph.Series) != 1 {
				t.Fatalf("unexpected series count: %d", len(graph.Series))
			}
			if graph.Series[0].Name != "temperature_control_current" {
				t.Fatalf("unexpected series name: %s", graph.Series[0].Name)
			}
			if len(graph.Series[0].Points) != 2 {
				t.Fatalf("unexpected point count: %d", len(graph.Series[0].Points))
			}
			if graph.Series[0].Points[1].Value != 20.1 {
				t.Fatalf("unexpected point value: %f", graph.Series[0].Points[1].Value)
			}
		})
	}
}

func TestQueryPrometheusGraphFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"error","error":"bad query"}`))
	}))
	defer server.Close()

	h := &httpd{
		httpClient:    server.Client(),
		prometheusURL: server.URL,
	}

	if _, err := h.queryPrometheusGraph("living_room", "24h", graphRanges["24h"]); err == nil {
		t.Fatal("expected prometheus error")
	}
}
