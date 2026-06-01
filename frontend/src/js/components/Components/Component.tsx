import type React from "react";

import {
  BinaryFan,
  BinaryLight,
  BinarySensor,
  BinaryTRV,
  Boiler,
} from "./BinaryDevices";
import { DeviceStatus, GenericSensor, PowerMeter, WifiSignal } from "./Sensors";
import { Switch, TemperatureSwitch, VirtualSwitch } from "./Switches";

import { ClimateSensor } from "./ClimateSensor";
import { EsphomeLight } from "./EsphomeLight";
import { HomedTemperature } from "./HomedTemperature";
import { RollerShutter } from "./RollerShutter";
import { HeaderCard } from "./common/HeaderCard";

interface ComponentProps {
  id: string;
  type: string;
  noCard?: boolean;
}

export const Component: React.FC<ComponentProps> = ({
  id,
  type,
  noCard = false,
}) => {
  let typedComponent: React.ReactNode;
  switch (type) {
    case "generic_sensor":
      typedComponent = <GenericSensor id={id} />;
      break;
    case "binary_sensor":
      typedComponent = <BinarySensor id={id} />;
      break;
    case "switch":
      typedComponent = <Switch id={id} />;
      break;
    case "homed_temperature":
      typedComponent = <HomedTemperature id={id} />;
      break;
    case "homed_temperature_switch":
      typedComponent = <TemperatureSwitch id={id} />;
      break;
    case "roller_shutter":
      typedComponent = <RollerShutter id={id} />;
      break;
    case "power_meter":
      typedComponent = <PowerMeter id={id} />;
      break;
    case "esphome_light":
      typedComponent = <EsphomeLight id={id} />;
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
    case "zigbee_climate_sensor":
      typedComponent = <ClimateSensor id={id} />;
      break;
    case "device_status":
      typedComponent = <DeviceStatus id={id} />;
      break;
    case "wifi_signal":
      typedComponent = <WifiSignal id={id} />;
      break;
    case "virtual_switch":
      typedComponent = <VirtualSwitch id={id} />;
      break;
    case "weather":
      typedComponent = <ClimateSensor id={id} />;
      break;
    default:
      return <h3>Unknown component type: {type}</h3>;
  }

  if (noCard) {
    return <>{typedComponent}</>;
  }

  return <HeaderCard id={id}>{typedComponent}</HeaderCard>;
};
