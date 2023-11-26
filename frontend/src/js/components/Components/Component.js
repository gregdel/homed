import React from "react";
import PropTypes from "prop-types";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
dayjs.extend(relativeTime);

import Icon from "@mdi/react";
import { mdiChartLine } from "@mdi/js";

import { Link } from "react-router-dom";

import { Switch } from "./Switch";
import { PowerMeter } from "./PowerMeter";
import { BinaryLight } from "./BinaryLight";
import { BinaryFan } from "./BinaryFan";
import { BinaryTRV } from "./BinaryTRV";

import { Boiler } from "./Boiler";
import { DeviceStatus } from "./DeviceStatus";
import { RollerShutter } from "./RollerShutter";
import { GenericSensor } from "./GenericSensor";
import { BinarySensor } from "./BinarySensor";
import { HomedTemperature } from "./HomedTemperature";
import { TemperatureSwitch } from "./TemperatureSwitch";
import { ZigbeeTRV } from "./ZigbeeTRV";
import { WifiSignal } from "./WifiSignal";
import { ZigbeeClimateSensor } from "./ZigbeeClimateSensor";
import { VirtualSwitch } from "./VirtualSwitch";

import { Card } from "antd";

export const Component = ({
  id,
  type,
  title,
  roomName,
  deviceName,
  online,
  friendlyName,
  updatedAt,
  graphURL,
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
    case "binary_fan":
      typedComponent = <BinaryFan id={id} />;
      break;
    case "binary_trv":
      typedComponent = <BinaryTRV id={id} />;
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
    case "homed_temperature_switch":
      typedComponent = <TemperatureSwitch id={id} />;
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
    case "virtual_switch":
      typedComponent = <VirtualSwitch id={id} />;
      break;
    default:
      typedComponent = <>Unhandled {type}</>;
      break;
  }

  if (noCard) {
    return typedComponent;
  }

  var extras = [];
  extras.push(
    <div key="status">
      {online === true ? dayjs(updatedAt).fromNow() : "offline"}
    </div>
  );

  if (graphURL && graphURL != "") {
    extras.push(
      <div key="graphIcon" style={{ marginRight: "-1em", marginLeft: "0.3em" }}>
        <Link to={`/components/${id}/graph`}>
          <Icon path={mdiChartLine} style={{ color: "#000000d9" }} size={0.8} />
        </Link>
      </div>
    );
  }

  const newTitle =
    friendlyName !== ""
      ? friendlyName
      : `${roomName} - ${type} - ${deviceName}`;
  return (
    <Card
      title={title ? title : newTitle}
      extra={<div style={{ display: "flex" }}>{extras}</div>}
    >
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
  noCard: PropTypes.bool.isRequired,
  graphURL: PropTypes.string,
};
