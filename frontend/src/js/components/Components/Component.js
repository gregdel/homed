import React from "react";
import PropTypes from "prop-types";
import moment from "moment";

import { Switch } from "./Switch";
import { PowerMeter } from "./PowerMeter";
import { BinaryLight } from "./BinaryLight";

import { Boiler } from "./Boiler";
import { DeviceStatus } from "./DeviceStatus";
import { RollerShutter } from "./RollerShutter";
import { GenericSensor } from "./GenericSensor";
import { BinarySensor } from "./BinarySensor";
import { HomedTemperature } from "./HomedTemperature";
import { ZigbeeTRV } from "./ZigbeeTRV";
import { WifiSignal } from "./WifiSignal";
import { ZigbeeClimateSensor } from "./ZigbeeClimateSensor";

import { Card } from "antd";

export const Component = ({
  id,
  type,
  title,
  extra,
  roomName,
  deviceName,
  online,
  friendlyName,
  updatedAt,
  noCard,
}) => {
  var typedComponent;
  switch (type) {
    case "generic_sensor":
      typedComponent = <GenericSensor id={id} />;
      break;
    case "binary_sensor":
      typedComponent = <BinarySensor id={id} />;
      break;
    case "device_status":
      typedComponent = <DeviceStatus id={id} />;
      break;
    case "wifi_signal":
      typedComponent = <WifiSignal id={id} />;
      break;
    case "boiler":
      typedComponent = <Boiler id={id} />;
      break;
    case "binary_light":
      typedComponent = <BinaryLight id={id} />;
      break;
    case "switch":
      typedComponent = <Switch id={id} />;
      break;
    case "roller_shutter":
      typedComponent = <RollerShutter id={id} />;
      break;
    case "homed_temperature":
      typedComponent = <HomedTemperature id={id} />;
      break;
    case "zigbee_climate_sensor":
      typedComponent = <ZigbeeClimateSensor id={id} />;
      break;
    case "power_meter":
      typedComponent = <PowerMeter id={id} />;
      break;
    case "zigbee_trv":
      typedComponent = <ZigbeeTRV id={id} />;
      break;
    default:
      typedComponent = <>Unhandled {type}</>;
      break;
  }

  if (noCard) {
    return typedComponent;
  }

  var status = "offline";
  if (online === true) {
    const m = moment(updatedAt, "YYYY-MM-DD HH:mm:ss Z");
    if (m.isValid()) {
      status = m.fromNow();
    }
  }

  const newTitle =
    friendlyName !== ""
      ? friendlyName
      : `${roomName} - ${type} - ${deviceName}`;
  return (
    <Card title={title ? title : newTitle} extra={extra ? extra : status}>
      {typedComponent}
    </Card>
  );
};
Component.defaultProps = {
  noCard: false,
};
Component.propTypes = {
  id: PropTypes.string.isRequired,
  type: PropTypes.string.isRequired,
  friendlyName: PropTypes.string,
  roomName: PropTypes.string,
  deviceName: PropTypes.string,
  online: PropTypes.bool,
  updatedAt: PropTypes.string,
  title: PropTypes.string,
  extra: PropTypes.any,
  noCard: PropTypes.bool.isRequired,
};
