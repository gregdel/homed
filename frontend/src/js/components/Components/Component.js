import React from "react";
import PropTypes from "prop-types";
import moment from "moment";

import { HomedTemperature } from "./HomedTemperature";
import { Temperature } from "./Temperature";
import { Humidity } from "./Humidity";
import { DeviceStatus } from "./DeviceStatus";
import { TuyaTRV } from "./TuyaTRV";
import { SaswellTRV } from "./SaswellTRV";
import { WifiSignal } from "./WifiSignal";
import { RTL433 } from "./RTL433";
import { Boiler } from "./Boiler";
import { TasmotaSwitch } from "./TasmotaSwitch";
import { ESPLight } from "./EspLight";
import { ESPCover } from "./EspCover";
import { XiaomiAqara } from "./XiaomiAqara";

import { Card } from "antd";

export const Component = ({
  id,
  type,
  title,
  extra,
  roomName,
  deviceName,
  friendlyName,
  updatedAt,
  noCard,
}) => {
  var typedComponent;
  switch (type) {
    case "temperature":
      typedComponent = <Temperature id={id} />;
      break;
    case "humidity":
      typedComponent = <Humidity id={id} />;
      break;
    case "device_status":
      typedComponent = <DeviceStatus id={id} />;
      break;
    case "tuya_trv":
      typedComponent = <TuyaTRV id={id} />;
      break;
    case "saswell_trv":
      typedComponent = <SaswellTRV id={id} />;
      break;
    case "wifi_signal":
      typedComponent = <WifiSignal id={id} />;
      break;
    case "rtl_433":
      typedComponent = <RTL433 id={id} />;
      break;
    case "boiler":
      typedComponent = <Boiler id={id} />;
      break;
    case "tasmota_switch":
      typedComponent = <TasmotaSwitch id={id} />;
      break;
    case "esphome_light":
      typedComponent = <ESPLight id={id} />;
      break;
    case "roller_shutter":
      typedComponent = <ESPCover id={id} />;
      break;
    case "homed_temperature":
      typedComponent = <HomedTemperature id={id} />;
      break;
    case "xiaomi_aqara":
      typedComponent = <XiaomiAqara id={id} />;
      break;
    default:
      typedComponent = <>Unhandled {type}</>;
      break;
  }

  if (noCard) {
    return typedComponent;
  }

  var prettyDate = "";
  const m = moment(updatedAt, "YYYY-MM-DD HH:mm:ss Z");
  if (m.isValid()) {
    prettyDate = m.fromNow();
  }

  const newTitle =
    friendlyName !== ""
      ? friendlyName
      : `${roomName} - ${type} - ${deviceName}`;
  return (
    <Card title={title ? title : newTitle} extra={extra ? extra : prettyDate}>
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
  updatedAt: PropTypes.string,
  title: PropTypes.string,
  extra: PropTypes.any,
  noCard: PropTypes.bool.isRequired,
};
