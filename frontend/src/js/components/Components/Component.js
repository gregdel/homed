import React from "react";
import PropTypes from "prop-types";
import moment from "moment";

import { Boiler } from "./Boiler";
import { DeviceStatus } from "./DeviceStatus";
import { ESPCover } from "./EspCover";
import { ESPLight } from "./EspLight";
import { ESPSwitch } from "./EspSwitch";
import { HomedTemperature } from "./HomedTemperature";
import { Humidity } from "./Humidity";
import { PowerMeter } from "./PowerMeter";
import { RTL433 } from "./RTL433";
import { SaswellTRV } from "./SaswellTRV";
import { Temperature } from "./Temperature";
import { TuyaTRV } from "./TuyaTRV";
import { WifiSignal } from "./WifiSignal";
import { XiaomiAqara } from "./XiaomiAqara";

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
    case "esphome_light":
      typedComponent = <ESPLight id={id} />;
      break;
    case "esphome_switch":
      typedComponent = <ESPSwitch id={id} />;
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
    case "power_meter":
      typedComponent = <PowerMeter id={id} />;
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
