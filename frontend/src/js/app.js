import React from "react";
import ReactDOM from "react-dom";

import { Provider } from "react-redux";
import store from "./store";

import { Layout } from "antd";
const { Content } = Layout;

import { WsHandler } from "./websocket";
import { AppMenu } from "./components/Menu";
import { Dashboard } from "./components/Dashboard";

import "antd/dist/antd.css";
import "../css/index.css";

const App = () => (
  <Provider store={store}>
    <WsHandler />
    <Layout>
      <AppMenu />
      <Layout>
        <Content className="homed-content">
          <Dashboard />
        </Content>
      </Layout>
    </Layout>
  </Provider>
);

ReactDOM.render(<App />, document.getElementById("app"));
