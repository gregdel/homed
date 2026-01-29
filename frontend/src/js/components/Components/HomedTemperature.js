import React, { useState, useEffect, useRef } from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

import { prettyName } from "../../utils";
import { relativeTime } from "../../utils/relativeTime";

import { Icon } from "../ui/Icon";
import { Card } from "../ui/Card";
import { Slider } from "../ui/Slider";

import { Link } from "../Navigation";

// Popover component using details/summary
const Popover = ({ trigger, title, children }) => {
  const detailsRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (e) => {
      if (detailsRef.current && !detailsRef.current.contains(e.target)) {
        detailsRef.current.open = false;
      }
    };
    document.addEventListener("click", handleClickOutside);
    return () => document.removeEventListener("click", handleClickOutside);
  }, []);

  return (
    <details ref={detailsRef} style={{ position: "relative" }}>
      <summary className="cursor-pointer" style={{ listStyle: "none" }}>
        {trigger}
      </summary>
      <div
        style={{
          position: "absolute",
          top: "100%",
          left: "50%",
          transform: "translateX(-50%)",
          background: "#fff",
          border: "1px solid #d9d9d9",
          borderRadius: "6px",
          padding: "0.5rem",
          zIndex: 10,
          boxShadow: "0 6px 16px rgba(0,0,0,0.08)",
          minWidth: "180px",
        }}
      >
        {title && (
          <div style={{ fontWeight: 500, marginBottom: "0.5rem" }}>{title}</div>
        )}
        {children}
      </div>
    </details>
  );
};

Popover.propTypes = {
  trigger: PropTypes.node.isRequired,
  title: PropTypes.string,
  children: PropTypes.node.isRequired,
};

export const HomedTemperature = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const {
    current,
    friendly_name: friendlyName,
    graph_url: graphURL,
    target,
    mode,
    device,
    heating,
    on,
    opportunistic,
    manual_target: manualTarget,
    manual_until: manualUntil,
  } = component.values;

  const [newTarget, setNewTarget] = useState(
    mode === "auto" ? target : manualTarget
  );

  useEffect(() => {
    setNewTarget(mode === "auto" ? target : manualTarget);
  }, [mode, target, manualTarget]);

  const [showTimePicker, setShowTimePicker] = useState(false);

  const sendChange = ({
    mode,
    target,
    duration = null,
    date = null,
    on = null,
  }) => {
    const data = {
      mode,
      manual_target: target,
      manual_until: date,
      manual_duration: duration,
      on,
    };
    setShowTimePicker(false);

    updateComponent(id, data);
  };

  const onChangeComplete = (value) => {
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
            duration: e.target.value + "h",
          });
        }}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            sendChange({
              mode: "duration",
              target: newTarget,
              duration: e.target.value + "h",
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

  const icon = opportunistic ? "leaf" : "radiator";

  const title = friendlyName !== "" ? friendlyName : prettyName(device.room);

  return (
    <Card
      title={title}
      style={{
        userSelect: "none",
      }}
      extra={
        <div className="flex items-center">
          {graphURL !== "" && (
            <>
              <Link
                to={`/components/${id}/graph`}
                style={{ color: "#000000d9" }}
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
                color: on ? "#000000" : "#00000040",
                transition: "color 0.3s ease-out 0s",
              }}
            />
          </div>
          <span className="divider-vertical" />
          <Link
            to={`/components/${id}/schedule`}
            style={{ color: "#000000d9" }}
          >
            <Icon name="calendarClock" size={1} />
          </Link>
        </div>
      }
      actions={actions}
    >
      <div style={{ fontSize: "4em" }}>
        <span>{current.toFixed(1)}°C</span>
      </div>

      {on && (
        <div>
          <div className="flex justify-between" style={{ fontSize: "1em" }}>
            <span>
              Target set to {mode === "auto" ? target : manualTarget}°C
            </span>
            <Icon
              name={icon}
              size={1}
              style={{
                color: heating ? "#ff00005e" : "#00000040",
              }}
            />
          </div>

          <div style={{ fontSize: "1em" }}>
            <HeatingMode mode={mode} date={manualUntil} />
          </div>

          <div style={{ position: "relative", marginTop: "1rem" }}>
            {target !== newTarget && (
              <div
                style={{
                  position: "absolute",
                  top: "-1.5rem",
                  left: `${((target - 8) / 17) * 100}%`,
                  transform: "translateX(-50%)",
                  cursor: "pointer",
                  fontSize: "0.85em",
                  color: "var(--color-primary)",
                }}
                onClick={() => onChangeComplete(target)}
              >
                {target}°C
              </div>
            )}
            <Slider
              min={8}
              max={25}
              step={0.5}
              value={newTarget}
              onChange={(value) => {
                setNewTarget(value);
              }}
              onChangeComplete={onChangeComplete}
            />
          </div>
        </div>
      )}
      {!on && (
        <div
          style={{
            fontSize: "2em",
            marginTop: "0px",
            marginBotton: "0px",
          }}
        >
          OFF
        </div>
      )}
    </Card>
  );
};
HomedTemperature.propTypes = {
  id: PropTypes.string.isRequired,
};

const HeatingMode = ({ mode, date }) => {
  const prettyDate = date === null ? "" : relativeTime(date);

  switch (mode) {
    case "auto":
      return <>Mode auto</>;
    case "fixed":
      return <>Manual until stopped</>;
    case "until_date":
    case "duration":
      return <>Back to auto {prettyDate}</>;
    case "until_next_change":
      if (date) {
        return <>Back to auto {prettyDate}</>;
      } else {
        return <>Manual until next programmed change</>;
      }
    default:
      return null;
  }
};
HeatingMode.propTypes = {
  mode: PropTypes.string.isRequired,
  date: PropTypes.string,
};
