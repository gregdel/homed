import React, { useState } from "react";
import PropTypes from "prop-types";
import { useDispatch, useSelector } from "react-redux";

import { Row, Col, Slider, Progress } from "antd";

import Icon from "@mdi/react";
import { mdiBattery, mdiBattery10 } from "@mdi/js";

import { componentUpdate } from "../../actions/homedComponents";

export const TuyaTRV = ({ uuid }) => {
  const dispatch = useDispatch();
  const {
    local_temperature: localTemperature,
    battery_low: batteryLow,
    current_heating_setpoint: currentHeatingSetpoint,
    position,
  } = useSelector((state) => state.stuff.components.get(uuid).values);

  const [target, setTarget] = useState(currentHeatingSetpoint);

  var marks = {};
  marks[localTemperature] = localTemperature + "°C";

  const onChange = (value) => {
    setTarget(value);
  };

  const onAfterChange = (value) => {
    dispatch(componentUpdate(uuid, value));
  };

  return (
    <>
      <Row justify="space-between" align="middle">
        <Col>
          <Icon
            path={batteryLow ? mdiBattery10 : mdiBattery}
            size={1}
            rotate={90}
          />
        </Col>
        <Progress percent={position} steps={5} />
      </Row>

      <div style={{ fontSize: "4em" }}>
        <span>{localTemperature}°C</span>
      </div>

      <div style={{ fontSize: "1em" }}>
        <span>Heating to {target}°C</span>
      </div>

      <Slider
        min={5}
        max={35}
        step={0.5}
        marks={marks}
        value={typeof target === "number" ? target : 0}
        onChange={onChange}
        onAfterChange={onAfterChange}
      />
    </>
  );
};

TuyaTRV.propTypes = {
  uuid: PropTypes.string.isRequired,
};
