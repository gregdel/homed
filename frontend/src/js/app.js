import React from "react";
import ReactDOM from "react-dom";
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
import { Dashboard } from "./components/TemperatureControl/Dashboard";
import { Components } from "./components/Components/Components";
import { Schedule } from "./components/Schedule/Schedule";

import store, { history } from "./store";

import "antd/dist/antd.css";
import "../css/index.css";

const App = () => (
  <Provider store={store}>
    <Router history={history}>
      <DataFetcher>
        <Layout style={{ minHeight: "100vh" }}>
          <AppMenu />
          <Layout>
            <Content style={{ padding: "1em" }}>
              <Switch>
                <Route path="/all" exact component={Components} />
                <Route path="/temperature" exact component={Dashboard} />
                <Route path="/tuya" exact>
                  <Components typesFilter={["tuya_trv"]} />
                </Route>
                <Route path="/lights" exact>
                  <Components typesFilter={["esphome_light"]} />
                </Route>
                <Route path="/switches" exact>
                  <Components typesFilter={["tasmota_switch"]} />
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

ReactDOM.render(<App />, document.getElementById("app"));
