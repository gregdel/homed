import React from "react";
import { Layout, Menu, Typography } from "antd";
import { Link } from "react-router-dom";

import Icon from "@mdi/react";
import {
  mdiMonitorDashboard,
  mdiHomeThermometer,
  mdiAvTimer,
  mdiThermometerLines,
} from "@mdi/js";

const { Sider } = Layout;

export const AppMenu = () => (
  <Sider breakpoint="lg" collapsedWidth="0">
    <Menu theme="dark" mode="inline" style={{ paddingTop: "1em" }}>
      <Menu.Item key="1" icon={<Icon path={mdiMonitorDashboard} size={1} />}>
        <Link to="/all" component={Typography.Link}>
          All
        </Link>
      </Menu.Item>
      <Menu.Item key="2" icon={<Icon path={mdiHomeThermometer} size={1} />}>
        <Link to="/temperature" component={Typography.Link}>
          Temperature
        </Link>
      </Menu.Item>
      <Menu.Item key="3" icon={<Icon path={mdiAvTimer} size={1} />}>
        <Link to="/schedule" component={Typography.Link}>
          Schedule
        </Link>
      </Menu.Item>
      <Menu.Item key="4" icon={<Icon path={mdiThermometerLines} size={1} />}>
        <Link to="/tuya" component={Typography.Link}>
          Tuya TRVs
        </Link>
      </Menu.Item>
    </Menu>
  </Sider>
);
