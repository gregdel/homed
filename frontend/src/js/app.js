import React from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import {
  HashRouter as Router,
  Route,
  Routes,
  Navigate,
} from "react-router-dom";

import { Layout } from "antd";
const { Content } = Layout;

import { AppMenu } from "./components/Menu";
import { DataFetcher } from "./components/DataFetcher";
import { Notifications } from "./components/Notifications";
import { Dashboard } from "./components/TemperatureControl/Dashboard";
import { Components } from "./components/Components/Components";
import { Schedule } from "./components/Schedule/Schedule";

import store from "./store";

import "@ant-design/cssinjs";
import "../assets/app.css";

const App = () => (
  <Provider store={store}>
    <Router>
      <DataFetcher>
        <Layout style={{ minHeight: "100vh" }}>
          <AppMenu />
          <Layout>
            <Content style={{ padding: "1em" }}>
              <Notifications />
              <Routes>
                <Route path="/all" exact element={<Components />} />
                <Route path="/temperature" exact element={<Dashboard />} />
                <Route
                  path="/trv"
                  exact
                  element={
                    <Components typesFilter={["zigbee_trv", "binary_trv"]} />
                  }
                />
                <Route
                  path="/lights"
                  exact
                  element={<Components typesFilter={["binary_light"]} />}
                />
                <Route
                  path="/power"
                  exact
                  element={<Components typesFilter={["power_meter"]} />}
                />
                <Route
                  path="/switches"
                  exact
                  element={<Components typesFilter={["switch"]} />}
                />
                <Route
                  path="/fans"
                  exact
                  element={<Components typesFilter={["binary_fan"]} />}
                />
                <Route
                  path="/shutters"
                  exact
                  element={<Components typesFilter={["roller_shutter"]} />}
                />
                <Route
                  path="/sensors"
                  exact
                  element={
                    <Components
                      typesFilter={["generic_sensor", "binary_sensor"]}
                    />
                  }
                />
                <Route
                  path="/climate_sensors"
                  exact
                  element={
                    <Components typesFilter={["zigbee_climate_sensor"]} />
                  }
                />
                <Route
                  path="/components/:componentId/schedule"
                  exact
                  element={<Schedule />}
                />
                <Route path="*" element={<Navigate to="/temperature" />} />
              </Routes>
            </Content>
          </Layout>
        </Layout>
      </DataFetcher>
    </Router>
  </Provider>
);

const container = document.getElementById("app");
const root = createRoot(container);
root.render(<App />);
