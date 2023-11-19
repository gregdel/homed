import React, { useState, useEffect } from "react";
import { Layout, Menu, Typography } from "antd";
import { Link, useLocation } from "react-router-dom";

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
} from "@mdi/js";

const { Sider } = Layout;

export const AppMenu = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [collapsedType, setCollapsedType] = useState(null);
  const [selectedKeys, setSelectedKeys] = useState([]);
  const location = useLocation();

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

  const link = (title, path) => (
    <Link to={path} component={Typography.Link}>
      {title}
    </Link>
  );

  const icon = (path) => <Icon path={path} size={1} />;

  const items = [
    {
      icon: icon(mdiHomeThermometer),
      key: "temperature",
      label: link("Temperature", "/temperature"),
    },
    {
      icon: icon(mdiRadiator),
      label: link("Thermostatic valves", "/trv"),
      key: "trv",
    },
    {
      icon: icon(mdiThermometer),
      label: link("Temperature sensors", "/climate_sensors"),
      key: "climate_sensors",
    },
    {
      icon: icon(mdiLightbulb),
      label: link("Lights", "/lights"),
      key: "lights",
    },
    {
      icon: icon(mdiLightningBolt),
      label: link("Power consumption", "/power"),
      key: "power_consumption",
    },
    {
      icon: icon(mdiPower),
      label: link("Switches", "/switches"),
      key: "switches",
    },
    {
      icon: icon(mdiWindowShutter),
      label: link("Roller shutters", "/shutters"),
      key: "shutters",
    },
    {
      icon: icon(mdiRuler),
      label: link("Sensors", "/sensors"),
      key: "sensors",
    },
  ];

  useEffect(() => {
    let keys = [];
    items.map((entry) => {
      if (location.pathname === entry.label.props.to) {
        keys.push(entry.key);
      }
    });
    setSelectedKeys(keys);
  }, [location]);

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
        selectedKeys={selectedKeys}
        items={items}
      />
    </Sider>
  );
};
