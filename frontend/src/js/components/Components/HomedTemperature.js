import React, { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { useSelector, useDispatch } from "react-redux";
import { prettyName } from "../../utils";
import moment from "moment";

import { componentUpdate } from "../../actions/components";

import Icon from "@mdi/react";
import {
  mdiInfinity,
  mdiCalendar,
  mdiTimerOutline,
  mdiAutorenew,
  mdiAvTimer,
} from "@mdi/js";

import { Card, Slider, Popover, DatePicker, Input } from "antd";
import { Link } from "react-router-dom";

export const HomedTemperature = ({ id }) => {
  const dispatch = useDispatch();
  const {
    values: {
      current,
      target,
      mode,
      device,
      manual_target: manualTarget,
      manual_until: manualUntil,
    },
  } = useSelector((state) => state.components.components.get(id));

  const [newTarget, setNewTarget] = useState(
    mode === "auto" ? target : manualTarget
  );

  useEffect(() => {
    setNewTarget(mode === "auto" ? target : manualTarget);
  }, [mode]);

  const [showTimePicker, setShowTimePicker] = useState(false);

  const sendChange = ({ mode, target, duration = null, date = null }) => {
    const data = {
      mode,
      manual_target: target,
      manual_until: date,
      manual_duration: duration,
    };
    setShowTimePicker(false);

    dispatch(componentUpdate(id, data));
  };

  const onAfterChange = (value) => {
    if (value === target) {
      sendChange({ mode: "auto", target: value });
    } else {
      setShowTimePicker(true);
    }
  };

  var marks = {};
  if (target !== newTarget) {
    marks[target] = target + "°C";
  }

  const ActionInfinity = (
    <Icon
      path={mdiInfinity}
      size={1}
      onClick={() => {
        sendChange({ mode: "fixed", target: newTarget });
      }}
    />
  );

  const ActionCalendar = (
    <Popover
      title="Pick a date"
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
      <Icon path={mdiCalendar} size={1} />
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
      content={
        <Input
          type="number"
          onBlur={onActionTimerEvent}
          onPressEnter={onActionTimerEvent}
        />
      }
    >
      <Icon path={mdiTimerOutline} size={1} />
    </Popover>
  );

  const ActionAuto = (
    <Icon
      path={mdiAutorenew}
      size={1}
      onClick={() => {
        sendChange({ mode: "until_next_change", target: newTarget });
      }}
    />
  );

  let actions = showTimePicker
    ? [ActionTimer, ActionCalendar, ActionAuto, ActionInfinity]
    : [];

  return (
    <Card
      title={prettyName(device.room)}
      extra={
        <Link to={`/components/${id}/schedule`} style={{ color: "#000000d9" }}>
          <Icon path={mdiAvTimer} size={1} />
        </Link>
      }
      actions={actions}
    >
      <div style={{ fontSize: "4em" }}>
        <span>{current.toFixed(1)}°C</span>
      </div>

      <div style={{ fontSize: "1em" }}>
        <span>Heating to {mode === "auto" ? target : manualTarget}°C</span>
      </div>

      <div style={{ fontSize: "1em" }}>
        <HeatingMode mode={mode} date={manualUntil} />
      </div>

      <Slider
        min={5}
        max={35}
        step={0.5}
        marks={marks}
        value={newTarget}
        onChange={(value) => {
          setNewTarget(value);
        }}
        onAfterChange={onAfterChange}
      />
    </Card>
  );
};
HomedTemperature.propTypes = {
  id: PropTypes.string.isRequired,
};

const HeatingMode = ({ mode, date }) => {
  const m = moment(date, "YYYY-MM-DD HH:mm:ss Z");
  const prettyDate = m.isValid() ? m.fromNow() : date;

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
