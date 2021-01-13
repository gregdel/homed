import React, { useEffect } from "react";
import ReactDOM from "react-dom";
import PropTypes from "prop-types";
import { Provider, useDispatch } from "react-redux";
import {
  HashRouter as Router,
  Route,
  Switch,
  Redirect,
} from "react-router-dom";

import { Layout } from "antd";
const { Content } = Layout;

import { WsHandler } from "./websocket";
import { AppMenu } from "./components/Menu";
import { Dashboard } from "./components/TemperatureControl/Dashboard";
import { HomedComponents } from "./components/HomedComponents/Components";

import store, { history } from "./store";
import { fetchStuff } from "./actions/stuff";

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
                <Route path="/all" exact component={HomedComponents} />
                <Route path="/temperature" exact component={Dashboard} />
                <Route render={() => <Redirect to="/temperature" />} />
              </Switch>
            </Content>
          </Layout>
        </Layout>
      </DataFetcher>
    </Router>
  </Provider>
);

const DataFetcher = ({ children }) => {
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(fetchStuff());
  }, [dispatch]);

  return (
    <>
      <WsHandler />
      {children}
    </>
  );
};
DataFetcher.propTypes = {
  children: PropTypes.object,
};

ReactDOM.render(<App />, document.getElementById("app"));
