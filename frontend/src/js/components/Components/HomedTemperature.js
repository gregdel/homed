import React, { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

import { prettyName } from "../../utils";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
dayjs.extend(relativeTime);

import Icon from "@mdi/react";
import {
  mdiInfinity,
  mdiCalendar,
  mdiTimerOutline,
  mdiAutorenew,
  mdiCalendarClock,
  mdiChartLine,
  mdiThermometer,
  mdiThermometerOff,
  mdiRadiator,
  mdiLeaf,
} from "@mdi/js";

import { Card, Slider, Popover, DatePicker, Input, Divider } from "antd";
import { Link } from "../Navigation";

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
  }, [mode]);

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

  var marks = {};
  if (target !== newTarget) {
    marks[target] = (
      <span onClick={() => onChangeComplete(target)}>{target}°C</span>
    );
  }

  const ActionInfinity = (
    <div
      onClick={() => {
        sendChange({ mode: "fixed", target: newTarget });
      }}
    >
      <Icon path={mdiInfinity} size={1} />
    </div>
  );

  const ActionCalendar = (
    <Popover
      title="Pick a date"
      trigger={["hover", "click"]}
      content={
        <DatePicker
          showTime
          onOk={(time) => {
            sendChange({
              mode: "until_date",
              target: newTarget,
              date: time.format(),
            });
          }}
        />
      }
    >
      <div>
        <Icon path={mdiCalendar} size={1} />
      </div>
    </Popover>
  );

  const onActionTimerEvent = (e) => {
    sendChange({
      mode: "duration",
      target: newTarget,
      duration: e.target.value + "h",
    });
  };

  const ActionTimer = (
    <Popover
      title="Select a duration (in hours)"
      trigger={["hover", "click"]}
      content={
        <Input
          type="number"
          onBlur={onActionTimerEvent}
          onPressEnter={onActionTimerEvent}
        />
      }
    >
      <div>
        <Icon path={mdiTimerOutline} size={1} />
      </div>
    </Popover>
  );

  const ActionUntilNext = (
    <div
      onClick={() => {
        sendChange({ mode: "until_next_change", target: newTarget });
      }}
    >
      <Icon path={mdiAutorenew} size={1} />
    </div>
  );

  const actions =
    on && showTimePicker
      ? [ActionTimer, ActionCalendar, ActionUntilNext, ActionInfinity]
      : [];

  const icon = opportunistic ? mdiLeaf : mdiRadiator;

  const title = friendlyName !== "" ? friendlyName : prettyName(device.room);

  return (
    <Card
      title={title}
      style={{
        userSelect: "none",
      }}
      extra={
        <div style={{ display: "flex" }}>
          {graphURL !== "" && (
            <>
              <Link
                to={`/components/${id}/graph`}
                style={{ color: "#000000d9" }}
              >
                <Icon path={mdiChartLine} size={1} />
              </Link>
              <Divider type="vertical" style={{ height: "1.8rem" }} />
            </>
          )}
          <div
            onClick={() => {
              sendChange({ mode: "on_off", on: !on });
            }}
          >
            <Icon
              path={on ? mdiThermometer : mdiThermometerOff}
              size={1}
              style={{
                cursor: "pointer",
                color: on ? "#000000" : "#00000040",
                transition: "color 0.3s ease-out 0s",
              }}
            />
          </div>
          <Divider type="vertical" style={{ height: "1.8rem" }} />
          <Link
            to={`/components/${id}/schedule`}
            style={{ color: "#000000d9" }}
          >
            <Icon path={mdiCalendarClock} size={1} />
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
          <div
            style={{
              fontSize: "1em",
              display: "flex",
              justifyContent: "space-between",
            }}
          >
            <span>
              Target set to {mode === "auto" ? target : manualTarget}°C
            </span>
            <Icon
              path={icon}
              size={1}
              style={{
                color: heating ? "#ff00005e" : "#00000040",
              }}
            />
          </div>

          <div style={{ fontSize: "1em" }}>
            <HeatingMode mode={mode} date={manualUntil} />
          </div>

          <Slider
            min={8}
            max={25}
            step={0.5}
            marks={marks}
            value={newTarget}
            onChange={(value) => {
              setNewTarget(value);
            }}
            onChangeComplete={onChangeComplete}
          />
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
  const prettyDate = date === null ? "" : dayjs(date).fromNow();

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
