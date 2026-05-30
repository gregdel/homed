import type React from "react";
import { useEffect, useRef, useState } from "react";
import { useComponents } from "../ComponentsContext";

import { prettyName } from "../../utils";
import { relativeTime } from "../../utils/relativeTime";

import { Card } from "../ui/Card";
import { Icon } from "../ui/Icon";
import { Slider } from "../ui/Slider";

import type { TemperatureMode, TemperatureValues } from "../../types";
import { Link } from "../Navigation";

interface PopoverProps {
  trigger: React.ReactNode;
  title?: string;
  children: React.ReactNode;
}

const Popover: React.FC<PopoverProps> = ({ trigger, title, children }) => {
  const detailsRef = useRef<HTMLDetailsElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (
        detailsRef.current &&
        !detailsRef.current.contains(e.target as Node)
      ) {
        detailsRef.current.open = false;
      }
    };
    document.addEventListener("click", handleClickOutside);
    return () => document.removeEventListener("click", handleClickOutside);
  }, []);

  return (
    <details ref={detailsRef} className="popover">
      <summary className="cursor-pointer" style={{ listStyle: "none" }}>
        {trigger}
      </summary>
      <div className="popover-content">
        {title && (
          <div style={{ fontWeight: 500, marginBottom: "0.5rem" }}>{title}</div>
        )}
        {children}
      </div>
    </details>
  );
};

interface HomedTemperatureProps {
  id: string;
}

export const HomedTemperature: React.FC<HomedTemperatureProps> = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const values = component.values as TemperatureValues;
  const {
    current,
    friendly_name: friendlyName,
    target,
    mode,
    device,
    heating,
    on,
    opportunistic,
    manual_target: manualTarget,
    manual_until: manualUntil,
  } = values;

  const [newTarget, setNewTarget] = useState(
    mode === "auto" ? target : (manualTarget ?? target),
  );

  useEffect(() => {
    setNewTarget(mode === "auto" ? target : (manualTarget ?? target));
  }, [mode, target, manualTarget]);

  const [showTimePicker, setShowTimePicker] = useState(false);

  const sendChange = ({
    mode,
    target,
    duration = null,
    date = null,
    on = null,
  }: {
    mode: TemperatureMode;
    target?: number;
    duration?: string | null;
    date?: string | null;
    on?: boolean | null;
  }) => {
    const data = {
      mode,
      manual_target: target,
      manual_until: date,
      manual_duration: duration,
      on,
    };
    setShowTimePicker(false);

    void updateComponent(id, data);
  };

  const onChangeComplete = (value: number) => {
    if (value === target) {
      sendChange({ mode: "auto", target: value });
    } else {
      setShowTimePicker(true);
    }
  };

  const ActionInfinity = (
    <div
      onClick={() => {
        sendChange({ mode: "fixed", target: newTarget });
      }}
    >
      <Icon name="infinity" size={1} />
    </div>
  );

  const ActionCalendar = (
    <Popover title="Pick a date" trigger={<Icon name="calendar" size={1} />}>
      <input
        type="datetime-local"
        className="input"
        onChange={(e) => {
          const date = new Date(e.target.value);
          sendChange({
            mode: "until_date",
            target: newTarget,
            date: date.toISOString(),
          });
        }}
      />
    </Popover>
  );

  const ActionTimer = (
    <Popover
      title="Duration (hours)"
      trigger={<Icon name="timerOutline" size={1} />}
    >
      <input
        type="number"
        className="input"
        onBlur={(e) => {
          sendChange({
            mode: "duration",
            target: newTarget,
            duration: `${e.target.value}h`,
          });
        }}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            sendChange({
              mode: "duration",
              target: newTarget,
              duration: `${(e.target as HTMLInputElement).value}h`,
            });
          }
        }}
      />
    </Popover>
  );

  const ActionUntilNext = (
    <div
      onClick={() => {
        sendChange({ mode: "until_next_change", target: newTarget });
      }}
    >
      <Icon name="autorenew" size={1} />
    </div>
  );

  const actions =
    on && showTimePicker
      ? [ActionTimer, ActionCalendar, ActionUntilNext, ActionInfinity]
      : [];

  const heatingIcon = opportunistic ? "leaf" : "radiator";
  const title = friendlyName !== "" ? friendlyName : prettyName(device.room);
  const effectiveTarget = mode === "auto" ? target : (manualTarget ?? target);

  return (
    <Card
      title={title}
      style={{ userSelect: "none" }}
      extra={
        <div className="flex items-center">
          {component.has_graph && (
            <>
              <Link
                to={`/components/${id}/graph`}
                style={{ color: "var(--color-text)" }}
              >
                <Icon name="chartLine" size={1} />
              </Link>
              <span className="divider-vertical" />
            </>
          )}
          <div
            onClick={() => {
              sendChange({ mode: "on_off", on: !on });
            }}
          >
            <Icon
              name={on ? "thermometer" : "thermometerOff"}
              size={1}
              style={{
                cursor: "pointer",
                color: on ? "var(--color-text)" : "var(--color-icon-inactive)",
                transition: "color 0.3s ease-out 0s",
              }}
            />
          </div>
          <span className="divider-vertical" />
          <Link
            to={`/components/${id}/schedule`}
            style={{ color: "var(--color-text)" }}
          >
            <Icon name="calendarClock" size={1} />
          </Link>
        </div>
      }
      actions={actions}
    >
      {on ? (
        <div className="temp-card-body">
          {/* Temperature display */}
          <div className="temp-row">
            <div className="temp-main">
              <div className="temp-current">{current.toFixed(1)}°</div>
              <div
                className={`temp-target-label ${mode === "auto" ? "temp-target-auto" : "temp-target-override"}`}
              >
                Target {effectiveTarget}°
              </div>
            </div>
            <div className="temp-heating">
              <Icon
                name={heatingIcon}
                size={1.5}
                style={{
                  color: heating
                    ? "var(--color-heating)"
                    : "var(--color-icon-inactive)",
                  transition: "color 0.3s ease-out",
                }}
              />
            </div>
          </div>

          {/* Mode indicator - clickable only when manual (to return to auto) */}
          <div className="temp-mode">
            <HeatingMode
              mode={mode}
              date={manualUntil}
              {...(mode !== "auto" && {
                onToggle: () => sendChange({ mode: "auto", target }),
              })}
            />
          </div>

          {/* Slider with labels */}
          <div className="temp-slider-container">
            {target !== newTarget && (
              <>
                <div
                  className="temp-slider-mark-label"
                  style={{ left: `${((target - 8) / 17) * 100}%` }}
                  onClick={() => onChangeComplete(target)}
                >
                  {target}°
                </div>
                <div
                  className="slider-mark"
                  style={{ left: `${((target - 8) / 17) * 100}%` }}
                  onClick={() => onChangeComplete(target)}
                  title={`Reset to ${target}°C (auto)`}
                />
              </>
            )}
            <Slider
              min={8}
              max={25}
              step={0.5}
              value={newTarget}
              onChange={(value: number) => {
                setNewTarget(value);
              }}
              onChangeComplete={onChangeComplete}
              formatTooltip={(v) => `${v}°C`}
            />
            <div className="temp-slider-labels">
              <span>8°</span>
              <span>25°</span>
            </div>
          </div>
        </div>
      ) : (
        <div className="temp-card-body">
          <div className="temp-row">
            <div className="temp-current">{current.toFixed(1)}°</div>
            <div className="temp-off-badge">OFF</div>
          </div>
        </div>
      )}
    </Card>
  );
};

interface HeatingModeProps {
  mode: TemperatureMode;
  date: string | null | undefined;
  onToggle?: () => void;
}

const HeatingMode: React.FC<HeatingModeProps> = ({ mode, date, onToggle }) => {
  const prettyDate = date ? relativeTime(date) : "";
  const manualClass = "temp-mode-chip temp-mode-manual cursor-pointer";

  switch (mode) {
    case "auto":
      return <span className="temp-mode-chip">auto</span>;
    case "fixed":
      return (
        <span className={manualClass} onClick={onToggle}>
          Manual until stopped
        </span>
      );
    case "until_date":
    case "duration":
      return (
        <span className={manualClass} onClick={onToggle}>
          Back to auto {prettyDate}
        </span>
      );
    case "until_next_change":
      if (date) {
        return (
          <span className={manualClass} onClick={onToggle}>
            Back to auto {prettyDate}
          </span>
        );
      }
      return (
        <span className={manualClass} onClick={onToggle}>
          Manual until next change
        </span>
      );
    default:
      return null;
  }
};
