import React from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import {
  HashRouter as Router,
  Route,
  Switch,
  Redirect,
} from "react-router-dom";

import { Layout } from "antd";
const { Content } = Layout;

import { AppMenu } from "./components/Menu";
import { DataFetcher } from "./components/DataFetcher";
import { Notifications } from "./components/Notifications";
import { Dashboard } from "./components/TemperatureControl/Dashboard";
import { Components } from "./components/Components/Components";
import { Schedule } from "./components/Schedule/Schedule";

import store, { history } from "./store";

import "@ant-design/cssinjs";
import "../assets/app.css";

const App = () => (
  <Provider store={store}>
    <Router history={history}>
      <DataFetcher>
        <Layout style={{ minHeight: "100vh" }}>
          <AppMenu />
          <Layout>
            <Content style={{ padding: "1em" }}>
              <Notifications />
              <Switch>
                <Route path="/all" exact component={Components} />
                <Route path="/temperature" exact component={Dashboard} />
                <Route path="/trv" exact>
                  <Components typesFilter={["zigbee_trv", "binary_trv"]} />
                </Route>
                <Route path="/lights" exact>
                  <Components typesFilter={["binary_light"]} />
                </Route>
                <Route path="/power" exact>
                  <Components typesFilter={["power_meter"]} />
                </Route>
                <Route path="/switches" exact>
                  <Components typesFilter={["switch"]} />
                </Route>
                <Route path="/shutters" exact>
                  <Components typesFilter={["roller_shutter"]} />
                </Route>
                <Route path="/sensors" exact>
                  <Components
                    typesFilter={["generic_sensor", "binary_sensor"]}
                  />
                </Route>
                <Route path="/climate_sensors" exact>
                  <Components typesFilter={["zigbee_climate_sensor"]} />
                </Route>
                <Route
                  path="/components/:componentId/schedule"
                  exact
                  component={Schedule}
                />
                <Route render={() => <Redirect to="/temperature" />} />
              </Switch>
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
