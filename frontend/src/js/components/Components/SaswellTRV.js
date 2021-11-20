import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

import { Row, Col } from "antd";

import Icon from "@mdi/react";
import { mdiBattery, mdiBattery10 } from "@mdi/js";

export const SaswellTRV = ({ id }) => {
  const {
    local_temperature: localTemperature,
    local_temperature_calibration: calibration,
    battery_low: batteryLow,
    current_heating_setpoint: currentHeatingSetpoint,
    system_mode: mode,
  } = useSelector((state) => state.components.components.get(id).values);

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
      </Row>

      <div>
        <span style={{ fontSize: "4em" }}>{localTemperature}°C</span>
        <small style={{ marginLeft: "1em" }}>({calibration}°C)</small>
      </div>

      <div style={{ fontSize: "1em" }}>
        <span>Heating to {currentHeatingSetpoint}°C</span>
        <br />
        <span>
          Mode: <strong>{mode}</strong>
        </span>
        <br />
      </div>
    </>
  );
};

SaswellTRV.propTypes = {
  id: PropTypes.string.isRequired,
};
