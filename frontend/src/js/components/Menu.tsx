import React, { useState, useEffect } from "react";
import { useNav } from "./Navigation";

import { Icon, IconName } from "./ui/Icon";

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
];

export const AppMenu: React.FC = () => {
  const [collapsed, setCollapsed] = useState(window.innerWidth < 992);
  const { currentPath, navigate } = useNav();

  // Handle responsive collapse
  useEffect(() => {
    const handleResize = () => {
      setCollapsed(window.innerWidth < 992);
    };
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const handleItemClick = (path: string) => {
    navigate(path);
    // Collapse on mobile after navigation
    if (window.innerWidth < 992) {
      setCollapsed(true);
    }
  };

  return (
    <>
      <button
        className="sidebar-toggle"
        onClick={() => setCollapsed(!collapsed)}
        aria-label="Toggle menu"
      >
        <Icon name="menu" size={1} />
      </button>
      <nav className={`sidebar ${collapsed ? "collapsed" : ""}`}>
        <ul className="nav">
          {menuItems.map((item) => (
            <li
              key={item.key}
              className={`nav-item ${
                currentPath === item.path ? "active" : ""
              }`}
            >
              <button onClick={() => handleItemClick(item.path)}>
                <Icon name={item.icon} size={1} />
                <span>{item.label}</span>
              </button>
            </li>
          ))}
        </ul>
      </nav>
    </>
  );
};
