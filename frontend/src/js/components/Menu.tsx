import type React from "react";
import { useEffect, useState } from "react";
import { BottomNav } from "./BottomNav";
import { useNav } from "./Navigation";
import { Icon, type IconName } from "./ui/Icon";

interface MenuItem {
  icon: IconName;
  key: string;
  label: string;
  path: string;
}

const menuItems: MenuItem[] = [
  {
    icon: "homeThermometer",
    key: "temperature",
    label: "Temperature",
    path: "/temperature",
  },
  { icon: "radiator", key: "trv", label: "Thermostatic valves", path: "/trv" },
  {
    icon: "thermometer",
    key: "climate_sensors",
    label: "Temperature sensors",
    path: "/climate_sensors",
  },
  { icon: "lightbulb", key: "lights", label: "Lights", path: "/lights" },
  {
    icon: "lightningBolt",
    key: "power_consumption",
    label: "Power consumption",
    path: "/power",
  },
  { icon: "power", key: "switches", label: "Switches", path: "/switches" },
  { icon: "fan", key: "fans", label: "Fans", path: "/fans" },
  {
    icon: "windowShutter",
    key: "shutters",
    label: "Roller shutters",
    path: "/shutters",
  },
  { icon: "ruler", key: "sensors", label: "Sensors", path: "/sensors" },
  {
    icon: "automation",
    key: "automations",
    label: "Automations",
    path: "/automations",
  },
];

const MOBILE_BREAKPOINT = 992;

const Sidebar: React.FC = () => {
  const { currentPath, navigate } = useNav();

  return (
    <nav className="sidebar">
      <ul className="nav">
        {menuItems.map((item) => (
          <li
            key={item.key}
            className={`nav-item ${currentPath === item.path ? "active" : ""}`}
          >
            <button onClick={() => navigate(item.path)}>
              <Icon name={item.icon} size={1} />
              <span>{item.label}</span>
            </button>
          </li>
        ))}
      </ul>
    </nav>
  );
};

export const AppMenu: React.FC = () => {
  const [isMobile, setIsMobile] = useState(
    window.innerWidth < MOBILE_BREAKPOINT,
  );

  useEffect(() => {
    const handleResize = () => {
      setIsMobile(window.innerWidth < MOBILE_BREAKPOINT);
    };
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  if (isMobile) {
    return <BottomNav />;
  }

  return <Sidebar />;
};
