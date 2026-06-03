// Core component types

export interface Device {
  name: string;
  room: string;
  online: boolean;
}

export interface ComponentValues {
  id: string;
  friendly_name: string;
  hide: boolean;
  device: Device;
  updated_at: string; // ISO timestamp
  [key: string]: unknown; // Type-specific fields
}

export interface ComponentJSON {
  type: string;
  read_only: boolean;
  has_graph: boolean;
  values: ComponentValues;
}

// Specific component types

export interface SwitchValues extends ComponentValues {
  on: boolean;
}

export interface BinarySensorValues extends ComponentValues {
  on: boolean;
}

export interface GenericSensorValues extends ComponentValues {
  value: unknown;
}

export interface PowerMeterValues extends ComponentValues {
  power: number;
  energy: number;
}

export interface RollerShutterValues extends ComponentValues {
  value: number; // 0-100 percent open
}

export type TemperatureMode =
  | "auto"
  | "fixed"
  | "duration"
  | "until_date"
  | "until_next_change"
  | "on_off";

export interface TemperatureValues extends ComponentValues {
  current: number;
  target: number;
  mode: TemperatureMode;
  manual_target?: number;
  manual_until?: string | null; // ISO timestamp
  heating: boolean;
  opportunistic: boolean;
  on: boolean;
}

export interface WifiSignalValues extends ComponentValues {
  value: number; // signal strength
}

export interface DeviceStatusValues extends ComponentValues {
  value: string; // status text
}

export type EsphomeLightColorMode = "rgb" | "cwww";

export interface EsphomeLightValues extends ComponentValues {
  on: boolean;
  brightness: number;
  color_mode?: string;
  color_temp?: number;
  cold_white?: number;
  warm_white?: number;
  red?: number;
  green?: number;
  blue?: number;
}

export type ScriptStatus = "idle" | "running" | "success" | "error";

export interface ScriptValues extends ComponentValues {
  status: ScriptStatus;
  last_started_at: string | null;
  last_finished_at: string | null;
  last_duration_ms: number;
  last_exit_code: number | null;
  last_error: string;
}

// Component type registry
export const ComponentTypes = {
  GENERIC_SENSOR: "generic_sensor",
  BINARY_SENSOR: "binary_sensor",
  SWITCH: "switch",
  HOMED_TEMPERATURE: "homed_temperature",
  HOMED_TEMPERATURE_SWITCH: "homed_temperature_switch",
  ROLLER_SHUTTER: "roller_shutter",
  POWER_METER: "power_meter",
  ESPHOME_LIGHT: "esphome_light",
  BOILER: "boiler",
  BINARY_LIGHT: "binary_light",
  BINARY_FAN: "binary_fan",
  ZIGBEE_TRV: "zigbee_trv",
  BINARY_TRV: "binary_trv",
  ZIGBEE_CLIMATE_SENSOR: "zigbee_climate_sensor",
  DEVICE_STATUS: "device_status",
  WIFI_SIGNAL: "wifi_signal",
  VIRTUAL_SWITCH: "virtual_switch",
  WEATHER: "weather",
  SCRIPT: "script",
} as const;

export type ComponentType =
  (typeof ComponentTypes)[keyof typeof ComponentTypes];
