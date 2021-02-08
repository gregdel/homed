import React from "react";
import { Layout, Menu, Typography } from "antd";
import { Link } from "react-router-dom";

import Icon from "@mdi/react";
import {
  mdiHomeThermometer,
  mdiThermometerLines,
  mdiLightbulb,
  mdiPower,
} from "@mdi/js";

const { Sider } = Layout;

export const AppMenu = () => (
  <Sider breakpoint="lg" collapsedWidth="0">
    <Menu theme="dark" mode="inline" style={{ paddingTop: "1em" }}>
      <Menu.Item key="2" icon={<Icon path={mdiHomeThermometer} size={1} />}>
        <Link to="/temperature" component={Typography.Link}>
          Temperature
        </Link>
      </Menu.Item>
      <Menu.Item key="3" icon={<Icon path={mdiThermometerLines} size={1} />}>
        <Link to="/tuya" component={Typography.Link}>
          Termostatic valves
        </Link>
      </Menu.Item>
      <Menu.Item key="4" icon={<Icon path={mdiLightbulb} size={1} />}>
        <Link to="/lights" component={Typography.Link}>
          Lights
        </Link>
      </Menu.Item>
      <Menu.Item key="5" icon={<Icon path={mdiPower} size={1} />}>
        <Link to="/switches" component={Typography.Link}>
          Switches
        </Link>
      </Menu.Item>
    </Menu>
  </Sider>
);
