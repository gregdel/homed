import React from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { Layout } from "antd";
const { Content } = Layout;

import { AppMenu } from "./components/Menu";
import { NavigationProvider, useNav } from "./components/Navigation";
import { ComponentsProvider } from "./components/ComponentsContext";
import { NotificationsProvider } from "./components/NotificationsContext";

import { DataFetcher } from "./components/DataFetcher";
import { Notifications } from "./components/Notifications";
import { Dashboard } from "./components/TemperatureControl/Dashboard";
import { Components } from "./components/Components/Components";
import { Schedule } from "./components/Schedule/Schedule";
import { Graph } from "./components/Graph";

import store from "./store";

import "@ant-design/cssinjs";
import "../assets/app.css";

const AppContent = () => {
  const { currentPath, params } = useNav();

  // Route mapping function
  const getComponent = () => {
    // Extract component paths for better matching
    const path = currentPath.split("?")[0]; // Remove query parameters if any

    switch (path) {
      case "/all":
        return <Components />;
      case "/temperature":
        return <Dashboard />;
      case "/trv":
        return <Components typesFilter={["zigbee_trv", "binary_trv"]} />;
      case "/lights":
        return <Components typesFilter={["binary_light", "esphome_light"]} />;
      case "/power":
        return <Components typesFilter={["power_meter"]} />;
      case "/switches":
        return <Components typesFilter={["switch", "virtual_switch"]} />;
      case "/fans":
        return <Components typesFilter={["binary_fan"]} />;
      case "/shutters":
        return <Components typesFilter={["roller_shutter"]} />;
      case "/sensors":
        return <Components typesFilter={["generic_sensor", "binary_sensor"]} />;
      case "/climate_sensors":
        return <Components typesFilter={["zigbee_climate_sensor"]} />;
      case `/components/${params.componentId}/schedule`:
        return <Schedule />;
      case `/components/${params.componentId}/graph`:
        return <Graph />;
      default:
        // Redirect to temperature (replacing Navigate component)
        if (path !== "/temperature") {
          window.location.replace("#/temperature");
        }
        return <Dashboard />;
    }
  };

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <AppMenu />
      <Layout>
        <Content style={{ padding: "1em" }}>
          <Notifications />
          {getComponent()}
        </Content>
      </Layout>
    </Layout>
  );
};

const App = () => (
  <Provider store={store}>
    <NavigationProvider>
      <NotificationsProvider>
        <ComponentsProvider>
          <DataFetcher>
            <AppContent />
          </DataFetcher>
        </ComponentsProvider>
      </NotificationsProvider>
    </NavigationProvider>
  </Provider>
);

const container = document.getElementById("app");
const root = createRoot(container);
root.render(<App />);

export default App;
