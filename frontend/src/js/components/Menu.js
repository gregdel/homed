import React, { useState, useEffect } from "react";
import { Layout, Menu, Typography } from "antd";
import { Link, useLocation } from "react-router-dom";

import Icon from "@mdi/react";
import {
  mdiHomeThermometer,
  mdiThermometer,
  mdiThermometerLines,
  mdiLightbulb,
  mdiPower,
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

  const menu = [
    {
      path: "/temperature",
      title: "Temperature",
      icon: mdiHomeThermometer,
    },
    {
      icon: mdiThermometerLines,
      path: "/trv",
      title: "Thermostatic valves",
    },
    {
      icon: mdiThermometer,
      path: "/aqara",
      title: "Temperature sensors",
    },
    {
      icon: mdiLightbulb,
      path: "/lights",
      title: "Lights",
    },
    {
      icon: mdiPower,
      path: "/switches",
      title: "Switches",
    },
    {
      icon: mdiWindowShutter,
      path: "/shutters",
      title: "Roller shutters",
    },
  ];

  useEffect(() => {
    let keys = [];
    menu.map((entry, key) => {
      if (location.pathname === entry.path) {
        keys.push(key.toString());
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
      >
        {menu.map((entry, key) => {
          return (
            <Menu.Item
              key={key}
              icon={<Icon path={entry.icon} size={1} />}
              style={{ display: "flex", alignItems: "center" }}
            >
              <Link to={entry.path} component={Typography.Link}>
                {entry.title}
              </Link>
            </Menu.Item>
          );
        })}
      </Menu>
    </Sider>
  );
};
