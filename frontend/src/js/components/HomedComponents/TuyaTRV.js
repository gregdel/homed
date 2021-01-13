import React, { useState } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";

import { Row, Col, Typography, Slider, Progress } from "antd";
const { Title } = Typography;

import Icon from "@mdi/react";
import { mdiBattery, mdiBattery10 } from "@mdi/js";

import { componentUpdate } from "../../actions/homedComponents";

export const TuyaTRV = ({
  uuid,
  battery_low,
  current_heating_setpoint,
  local_temperature,
  position,
}) => {
  const dispatch = useDispatch();
  const [target, setTarget] = useState(current_heating_setpoint);

  var marks = {};
  marks[local_temperature] = local_temperature + "°C";

  const onChange = (value) => {
    setTarget(value);
  };

  const onAfterChange = (value) => {
    dispatch(componentUpdate(uuid, value));
  };

  return (
    <>
      <Row gutter={1}>
        <Col flex={1}>
          <Icon
            path={battery_low ? mdiBattery10 : mdiBattery}
            size={1}
            rotate={90}
          />
        </Col>
        <Col flex={2}>
          <Progress percent={position} steps={5} />
        </Col>
      </Row>
      <Title level={1}>{local_temperature}°C</Title>
      <Title level={3}>Set to {target}°C</Title>
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
  battery_low: PropTypes.bool.isRequired,
  current_heating_setpoint: PropTypes.number.isRequired,
  local_temperature: PropTypes.number.isRequired,
  position: PropTypes.number.isRequired,
};
