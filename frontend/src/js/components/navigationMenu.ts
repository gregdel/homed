import type { ComponentJSON } from "../types";
import { ComponentTypes } from "../types";
import type { IconName } from "./ui/Icon";

export interface NavigationMenuItem {
  componentTypes: string[];
  icon: IconName;
  key: string;
  label: string;
  mobileLabel: string;
  path: string;
}

export const navigationMenuItems: NavigationMenuItem[] = [
  {
    componentTypes: [ComponentTypes.HOMED_TEMPERATURE],
    icon: "homeThermometer",
    key: "temperature",
    label: "Temperature",
    mobileLabel: "Temperature",
    path: "/temperature",
  },
  {
    componentTypes: [ComponentTypes.BINARY_LIGHT, ComponentTypes.ESPHOME_LIGHT],
    icon: "lightbulb",
    key: "lights",
    label: "Lights",
    mobileLabel: "Lights",
    path: "/lights",
  },
  {
    componentTypes: [ComponentTypes.ROLLER_SHUTTER],
    icon: "windowShutter",
    key: "shutters",
    label: "Roller shutters",
    mobileLabel: "Shutters",
    path: "/shutters",
  },
  {
    componentTypes: [ComponentTypes.SWITCH, ComponentTypes.VIRTUAL_SWITCH],
    icon: "power",
    key: "switches",
    label: "Switches",
    mobileLabel: "Switches",
    path: "/switches",
  },
  {
    componentTypes: [ComponentTypes.BINARY_FAN],
    icon: "fan",
    key: "fans",
    label: "Fans",
    mobileLabel: "Fans",
    path: "/fans",
  },
  {
    componentTypes: [
      ComponentTypes.GENERIC_SENSOR,
      ComponentTypes.BINARY_SENSOR,
    ],
    icon: "ruler",
    key: "sensors",
    label: "Sensors",
    mobileLabel: "Sensors",
    path: "/sensors",
  },
  {
    componentTypes: [ComponentTypes.POWER_METER],
    icon: "lightningBolt",
    key: "power_consumption",
    label: "Power consumption",
    mobileLabel: "Power",
    path: "/power",
  },
  {
    componentTypes: [ComponentTypes.ZIGBEE_TRV, ComponentTypes.BINARY_TRV],
    icon: "radiator",
    key: "trv",
    label: "Thermostatic valves",
    mobileLabel: "TRVs",
    path: "/trv",
  },
  {
    componentTypes: [
      ComponentTypes.ZIGBEE_CLIMATE_SENSOR,
      ComponentTypes.WEATHER,
    ],
    icon: "thermometer",
    key: "climate_sensors",
    label: "Temperature sensors",
    mobileLabel: "Climate",
    path: "/climate_sensors",
  },
  {
    componentTypes: [ComponentTypes.SCRIPT],
    icon: "automation",
    key: "automations",
    label: "Automations",
    mobileLabel: "Automations",
    path: "/automations",
  },
];

const visibleComponents = (
  components: Record<string, ComponentJSON>,
): ComponentJSON[] =>
  Object.values(components).filter(
    (component) => component.values.hide !== true,
  );

export const getAvailableMenuItems = (
  components: Record<string, ComponentJSON>,
): NavigationMenuItem[] => {
  const visible = visibleComponents(components);

  return navigationMenuItems.filter((item) =>
    visible.some((component) => item.componentTypes.includes(component.type)),
  );
};

export const categoryPaths = new Set(
  navigationMenuItems.map((item) => item.path),
);
