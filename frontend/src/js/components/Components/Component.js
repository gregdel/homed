import React from "react";
import PropTypes from "prop-types";

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
import { WifiSignal } from "./WifiSignal";
import { ZigbeeClimateSensor } from "./ZigbeeClimateSensor";
import { VirtualSwitch } from "./VirtualSwitch";
import { HeaderCard } from "./common/HeaderCard";

export const Component = ({ id, type, noCard = false }) => {
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

  return <HeaderCard id={id}>{typedComponent}</HeaderCard>;
};
Component.propTypes = {
  id: PropTypes.string.isRequired,
  type: PropTypes.string.isRequired,
  noCard: PropTypes.bool,
};
