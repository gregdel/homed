package httpd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type graphRange struct {
	duration time.Duration
	step     time.Duration
}

var graphRanges = map[string]graphRange{
	"1h":  {duration: time.Hour, step: time.Minute},
	"6h":  {duration: 6 * time.Hour, step: time.Minute},
	"24h": {duration: 24 * time.Hour, step: 5 * time.Minute},
	"7d":  {duration: 7 * 24 * time.Hour, step: 30 * time.Minute},
}

type graphPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type graphSeries struct {
	Name   string       `json:"name"`
	Metric string       `json:"metric"`
	Points []graphPoint `json:"points"`
}

type graphResponse struct {
	ComponentID string        `json:"component_id"`
	Range       string        `json:"range"`
	StepSeconds int64         `json:"step_seconds"`
	Series      []graphSeries `json:"series"`
}

type prometheusQueryRangeResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Data   struct {
		Result []struct {
			Metric map[string]string `json:"metric"`
			Values [][]any           `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func (h *httpd) getComponentGraph(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if _, err := h.components.Get(id); err != nil {
		h.httpError(w, fmt.Sprintf("failed to get the component: %s", err))
		return
	}

	if !h.components.HasGraph(id) {
		h.httpError(w, "component has no graph metrics")
		return
	}

	if strings.TrimSpace(h.prometheusURL) == "" {
		h.httpError(w, "prometheus url is not configured")
		return
	}

	rangeName := r.URL.Query().Get("range")
	if rangeName == "" {
		rangeName = "24h"
	}
	graphRange, ok := graphRanges[rangeName]
	if !ok {
		h.httpError(w, "invalid graph range")
		return
	}

	data, err := h.queryPrometheusGraph(id, rangeName, graphRange)
	if err != nil {
		h.httpError(w, fmt.Sprintf("failed to query prometheus: %s", err))
		return
	}

	h.httpRenderJSON(w, data)
}

func (h *httpd) queryPrometheusGraph(componentID, rangeName string, graphRange graphRange) (*graphResponse, error) {
	promURL, err := url.Parse(strings.TrimRight(h.prometheusURL, "/") + "/api/v1/query_range")
	if err != nil {
		return nil, err
	}

	now := time.Now()
	query := promURL.Query()
	query.Set("query", fmt.Sprintf(`{__name__=~"homed_.+",id=%q}`, componentID))
	query.Set("start", strconv.FormatInt(now.Add(-graphRange.duration).Unix(), 10))
	query.Set("end", strconv.FormatInt(now.Unix(), 10))
	query.Set("step", strconv.FormatInt(int64(graphRange.step.Seconds()), 10))
	promURL.RawQuery = query.Encode()

	resp, err := h.httpClient.Get(promURL.String())
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected prometheus status: %s", resp.Status)
	}

	promResp := prometheusQueryRangeResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&promResp); err != nil {
		return nil, err
	}

	if promResp.Status != "success" {
		if promResp.Error != "" {
			return nil, errors.New(promResp.Error)
		}
		return nil, fmt.Errorf("query failed")
	}

	result := &graphResponse{
		ComponentID: componentID,
		Range:       rangeName,
		StepSeconds: int64(graphRange.step.Seconds()),
		Series:      []graphSeries{},
	}

	for _, promSeries := range promResp.Data.Result {
		metricName := promSeries.Metric["__name__"]
		series := graphSeries{
			Name:   graphSeriesName(promSeries.Metric),
			Metric: metricName,
			Points: []graphPoint{},
		}

		for _, value := range promSeries.Values {
			if len(value) != 2 {
				continue
			}

			timestamp, ok := value[0].(float64)
			if !ok {
				continue
			}

			valueString, ok := value[1].(string)
			if !ok {
				continue
			}

			valueFloat, err := strconv.ParseFloat(valueString, 64)
			if err != nil {
				continue
			}

			series.Points = append(series.Points, graphPoint{
				Timestamp: int64(timestamp),
				Value:     valueFloat,
			})
		}

		result.Series = append(result.Series, series)
	}

	sort.Slice(result.Series, func(i, j int) bool {
		return result.Series[i].Name < result.Series[j].Name
	})

	return result, nil
}

func graphSeriesName(metric map[string]string) string {
	name := strings.TrimPrefix(metric["__name__"], "homed_")
	componentType := metric["component_type"]
	if componentType == "" {
		return name
	}

	return strings.TrimPrefix(name, componentType+"_")
}
