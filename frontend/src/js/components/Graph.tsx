import type React from "react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { apiGet } from "../utils/api";
import { useComponents } from "./ComponentsContext";
import { useNav } from "./Navigation";
import { Card } from "./ui/Card";
import { Icon } from "./ui/Icon";
import { useResumeRefresh } from "./useResumeRefresh";

type GraphRange = "1h" | "6h" | "24h" | "7d";

interface GraphPoint {
  timestamp: number;
  value: number;
}

interface GraphSeries {
  name: string;
  metric: string;
  points: GraphPoint[];
}

interface GraphData {
  component_id: string;
  range: GraphRange;
  step_seconds: number;
  series: GraphSeries[];
}

const ranges: GraphRange[] = ["1h", "6h", "24h", "7d"];

const colors = ["#5e81ac", "#a3be8c", "#bf616a", "#d08770", "#8fbcbb"];

const preventChartMouseFocus = (
  _data: unknown,
  event: React.MouseEvent<SVGGraphicsElement>,
) => {
  event.preventDefault();
};

const tooltipStyle = {
  background: "var(--color-bg)",
  border: "1px solid var(--color-border)",
  borderRadius: "var(--radius)",
  boxShadow: "var(--shadow)",
  color: "var(--color-text)",
  padding: "0.5rem 0.75rem",
};

export const Graph: React.FC = () => {
  const { params, goBack } = useNav();
  const { getComponentById } = useComponents();
  const [range, setRange] = useState<GraphRange>("24h");
  const [data, setData] = useState<GraphData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const requestIDRef = useRef(0);
  const activeRequestRef = useRef<AbortController | null>(null);

  const componentId = params.componentId || "";
  const component = getComponentById(componentId);
  const titleText = component?.values.friendly_name || componentId;

  const fetchGraph = useCallback(async () => {
    if (!componentId) {
      return;
    }

    activeRequestRef.current?.abort();
    const controller = new AbortController();
    activeRequestRef.current = controller;

    const requestID = requestIDRef.current + 1;
    requestIDRef.current = requestID;

    try {
      setLoading(true);
      setError(null);
      const response = await apiGet<GraphData>(
        `/components/${componentId}/graph?range=${range}`,
        { cache: "no-store", signal: controller.signal },
      );
      if (requestIDRef.current === requestID) {
        setData(response.data);
      }
    } catch (err) {
      if (err instanceof DOMException && err.name === "AbortError") {
        return;
      }

      if (requestIDRef.current === requestID) {
        setError(err instanceof Error ? err.message : String(err));
      }
    } finally {
      if (activeRequestRef.current === controller) {
        activeRequestRef.current = null;
      }

      if (requestIDRef.current === requestID && !controller.signal.aborted) {
        setLoading(false);
      }
    }
  }, [componentId, range]);

  useResumeRefresh(fetchGraph);

  useEffect(() => {
    void fetchGraph();
    return () => {
      activeRequestRef.current?.abort();
      activeRequestRef.current = null;
      requestIDRef.current += 1;
    };
  }, [fetchGraph]);

  const series = data?.series ?? [];

  return (
    <Card
      title={
        <div className="flex items-center gap-sm">
          <button
            type="button"
            className="btn btn-icon"
            aria-label="Return to previous page"
            onClick={() => goBack("/temperature")}
          >
            <Icon name="arrowLeft" size={0.8} />
          </button>
          <span>{titleText}</span>
        </div>
      }
      extra={
        <div
          role="group"
          className="button-group button-group-sm graph-range-selector select-none"
          aria-label="Graph time range"
        >
          {ranges.map((item) => (
            <button
              key={item}
              type="button"
              className="button-group-item"
              aria-pressed={range === item}
              onClick={() => setRange(item)}
            >
              {item}
            </button>
          ))}
        </div>
      }
    >
      {loading && <div className="graph-state">Loading graph...</div>}
      {error && <div className="graph-state graph-error">{error}</div>}
      {!loading && !error && series.length === 0 && (
        <div className="graph-state">No graph data available.</div>
      )}
      {!loading && !error && series.length > 0 && (
        <div className="grid gap-md">
          {series.map((item, index) => (
            <LineChart
              key={`${item.metric}-${item.name}`}
              color={colors[index % colors.length] || "#5e81ac"}
              series={item}
            />
          ))}
        </div>
      )}
    </Card>
  );
};

const LineChart: React.FC<{ color: string; series: GraphSeries }> = ({
  color,
  series,
}) => {
  const chart = useMemo(() => {
    if (series.points.length === 0) {
      return null;
    }

    const points = [...series.points].sort((a, b) => a.timestamp - b.timestamp);
    const minTime = Math.min(...points.map((point) => point.timestamp));
    const maxTime = Math.max(...points.map((point) => point.timestamp));
    const minValue = Math.min(...points.map((point) => point.value));
    const maxValue = Math.max(...points.map((point) => point.value));
    const valuePadding = Math.max((maxValue - minValue) * 0.08, 0.5);

    return {
      points,
      minTime,
      maxTime,
      minValue,
      maxValue,
      valueDomain: [minValue - valuePadding, maxValue + valuePadding] as [
        number,
        number,
      ],
    };
  }, [series.points]);

  if (!chart) {
    return null;
  }

  const gradientId = `graph-gradient-${series.metric}-${series.name}`.replace(
    /[^a-zA-Z0-9_-]/g,
    "-",
  );

  return (
    <div className="graph-chart">
      <div className="graph-chart-title flex justify-between gap-md mb-sm">
        <span>{prettyMetricName(series.name)}</span>
        <span className="graph-chart-range">
          {formatValue(chart.minValue)} - {formatValue(chart.maxValue)}
        </span>
      </div>
      <div className="graph-chart-frame select-none">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart
            data={chart.points}
            margin={{ top: 16, right: 18, bottom: 8, left: 0 }}
            onMouseDown={preventChartMouseFocus}
          >
            <defs>
              <linearGradient id={gradientId} x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor={color} stopOpacity={0.42} />
                <stop offset="95%" stopColor={color} stopOpacity={0.04} />
              </linearGradient>
            </defs>
            <CartesianGrid
              stroke="var(--color-border)"
              strokeDasharray="3 3"
              strokeOpacity={0.7}
              vertical={false}
            />
            <XAxis
              axisLine={false}
              dataKey="timestamp"
              domain={[chart.minTime, chart.maxTime]}
              minTickGap={24}
              tick={{ fill: "var(--color-text-tertiary)", fontSize: 12 }}
              tickFormatter={formatShortTime}
              tickLine={false}
              type="number"
            />
            <YAxis
              axisLine={false}
              domain={chart.valueDomain}
              tick={{ fill: "var(--color-text-tertiary)", fontSize: 12 }}
              tickFormatter={formatValue}
              tickLine={false}
              width={48}
            />
            <Tooltip
              contentStyle={tooltipStyle}
              cursor={{ stroke: color, strokeOpacity: 0.18, strokeWidth: 2 }}
              formatter={(value) => [
                formatValue(Number(value)),
                prettyMetricName(series.name),
              ]}
              itemStyle={{ color: "var(--color-text)" }}
              labelStyle={{
                color: "var(--color-text-tertiary)",
                marginBottom: "0.25rem",
              }}
              labelFormatter={(label) => formatTime(Number(label))}
            />
            <Area
              activeDot={{ r: 4, stroke: "var(--color-bg)", strokeWidth: 2 }}
              dataKey="value"
              dot={false}
              fill={`url(#${gradientId})`}
              fillOpacity={1}
              isAnimationActive={false}
              name={prettyMetricName(series.name)}
              stroke={color}
              strokeWidth={2.5}
              type="monotone"
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};

const prettyMetricName = (name: string): string =>
  name
    .replace(/^homed_/, "")
    .replace(/_/g, " ")
    .replace(/\b\w/g, (letter) => letter.toUpperCase());

const formatValue = (value: number): string => {
  if (Math.abs(value) >= 100) {
    return value.toFixed(0);
  }
  return value.toFixed(1);
};

const formatTime = (timestamp: number): string =>
  new Date(timestamp * 1000).toLocaleString([], {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });

const formatShortTime = (timestamp: number): string =>
  new Date(timestamp * 1000).toLocaleString([], {
    month: "short",
    day: "numeric",
    hour: "2-digit",
  });
