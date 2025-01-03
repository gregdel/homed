import React, { useState } from "react";
import { Layout, Menu } from "antd";
import { useNav } from "./Navigation";

import Icon from "@mdi/react";
import {
  mdiHomeThermometer,
  mdiLightbulb,
  mdiLightningBolt,
  mdiPower,
  mdiRadiator,
  mdiRuler,
  mdiThermometer,
  mdiWindowShutter,
  mdiFan,
} from "@mdi/js";

const { Sider } = Layout;

export const AppMenu = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [collapsedType, setCollapsedType] = useState(null);
  const { currentPath, navigate } = useNav();

  const onCollapse = (collapsed, type) => {
    setCollapsed(collapsed);
    setCollapsedType(type);
  };

  const onClick = () => {
    if (collapsedType !== "clickTrigger") {
      return;
    }

    setCollapsed(true);
  };

  const icon = (path) => <Icon path={path} size={1} />;

  const items = [
    {
      icon: icon(mdiHomeThermometer),
      key: "temperature",
      label: "Temperature",
      onClick: () => navigate("/temperature"),
    },
    {
      icon: icon(mdiRadiator),
      key: "trv",
      label: "Thermostatic valves",
      onClick: () => navigate("/trv"),
    },
    {
      icon: icon(mdiThermometer),
      key: "climate_sensors",
      label: "Temperature sensors",
      onClick: () => navigate("/climate_sensors"),
    },
    {
      icon: icon(mdiLightbulb),
      key: "lights",
      label: "Lights",
      onClick: () => navigate("/lights"),
    },
    {
      icon: icon(mdiLightningBolt),
      key: "power_consumption",
      label: "Power consumption",
      onClick: () => navigate("/power"),
    },
    {
      icon: icon(mdiPower),
      key: "switches",
      label: "Switches",
      onClick: () => navigate("/switches"),
    },
    {
      icon: icon(mdiFan),
      key: "fans",
      label: "Fans",
      onClick: () => navigate("/fans"),
    },
    {
      icon: icon(mdiWindowShutter),
      key: "shutters",
      label: "Roller shutters",
      onClick: () => navigate("/shutters"),
    },
    {
      icon: icon(mdiRuler),
      key: "sensors",
      label: "Sensors",
      onClick: () => navigate("/sensors"),
    },
  ];

  const getSelectedKeys = () => {
    const item = items.find(
      (item) => currentPath === item.onClick.toString().match(/"([^"]+)"/)[1]
    );
    return item ? [item.key] : [];
  };

  return (
    <Sider
      breakpoint="lg"
      collapsedWidth="0"
      collapsed={collapsed}
      onCollapse={onCollapse}
    >
      <Menu
        theme="dark"
        mode="inline"
        style={{ paddingTop: "1em" }}
        onClick={onClick}
        selectedKeys={getSelectedKeys()}
        items={items}
      />
    </Sider>
  );
};
