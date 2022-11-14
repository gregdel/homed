import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

import { Row, Col, Progress } from "antd";

import Icon from "@mdi/react";
import { mdiBattery, mdiSync, mdiBattery10 } from "@mdi/js";

export const ZigbeeTRV = ({ id }) => {
  const {
    local_temperature: localTemperature,
    local_temperature_calibration: calibration,
    battery_low: batteryLow,
    current_heating_setpoint: currentHeatingSetpoint,
    calibration_request_time: calibrationRequestTime,
    system_mode: mode,
    force: force,
    position,
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
        {position >= 0 && <Progress percent={position} steps={5} />}
      </Row>

      <div>
        <span style={{ fontSize: "4em" }}>{localTemperature}°C</span>
        <small style={{ marginLeft: "1em" }}>
          ({calibration}°C
          {calibrationRequestTime !== null && (
            <>
              {" "}
              <Icon path={mdiSync} size={0.4} />
            </>
          )}
          )
        </small>
      </div>

      <div style={{ fontSize: "1em" }}>
        <span>Heating to {currentHeatingSetpoint}°C</span>
        <br />
        <span>
          Mode: <strong>{mode}</strong>
        </span>
        <br />
        {force !== "" && (
          <span>
            Force: <strong>{force}</strong>
          </span>
        )}
      </div>
    </>
  );
};

ZigbeeTRV.propTypes = {
  id: PropTypes.string.isRequired,
};
