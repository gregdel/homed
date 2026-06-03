import type React from "react";
import { useEffect, useMemo } from "react";
import { createRoot } from "react-dom/client";

import { AppMenu } from "./components/Menu";
import {
  ComponentsProvider,
  NavigationProvider,
  NotificationsProvider,
  useNav,
} from "./components/contexts";

import { Components } from "./components/Components/Components";
import { useComponents } from "./components/ComponentsContext";
import { DataFetcher } from "./components/DataFetcher";
import { Graph } from "./components/Graph";
import { Notifications } from "./components/Notifications";
import { Schedule } from "./components/Schedule/Schedule";
import { Dashboard } from "./components/TemperatureControl/Dashboard";
import {
  categoryPaths,
  getAvailableMenuItems,
} from "./components/navigationMenu";

import "../assets/app.css";

const AppContent: React.FC = () => {
  const { currentPath, params, navigate } = useNav();
  const { components, hasLoaded } = useComponents();
  const availableMenuItems = useMemo(
    () => getAvailableMenuItems(components),
    [components],
  );
  const firstAvailablePath = availableMenuItems[0]?.path;
  const path = currentPath.split("?")[0] ?? "";
  const isCategoryPath = categoryPaths.has(path);

  useEffect(() => {
    if (!hasLoaded || !firstAvailablePath) {
      return;
    }

    if (
      isCategoryPath &&
      !availableMenuItems.some((item) => item.path === path)
    ) {
      navigate(firstAvailablePath);
      return;
    }

    if (
      !isCategoryPath &&
      path !== "/all" &&
      params.componentId === undefined
    ) {
      navigate(firstAvailablePath);
    }
  }, [
    availableMenuItems,
    firstAvailablePath,
    hasLoaded,
    isCategoryPath,
    navigate,
    params.componentId,
    path,
  ]);

  // Route mapping function
  const getComponent = (): React.ReactNode => {
    if (hasLoaded && isCategoryPath && firstAvailablePath === undefined) {
      return (
        <div className="grid grid-cols-1 grid-cols-sm-2 grid-cols-lg-3">
          <div>No visible components.</div>
        </div>
      );
    }

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
      case "/automations":
        return <Components typesFilter={["script"]} />;
      case "/fans":
        return <Components typesFilter={["binary_fan"]} />;
      case "/shutters":
        return <Components typesFilter={["roller_shutter"]} />;
      case "/sensors":
        return <Components typesFilter={["generic_sensor", "binary_sensor"]} />;
      case "/climate_sensors":
        return (
          <Components typesFilter={["zigbee_climate_sensor", "weather"]} />
        );
      case `/components/${params.componentId}/schedule`:
        return <Schedule />;
      case `/components/${params.componentId}/graph`:
        return <Graph />;
      default:
        if (hasLoaded && firstAvailablePath === undefined) {
          return (
            <div className="grid grid-cols-1 grid-cols-sm-2 grid-cols-lg-3">
              <div>No visible components.</div>
            </div>
          );
        }
        return null;
    }
  };

  return (
    <div className="layout">
      <AppMenu />
      <main className="content">
        <Notifications />
        {getComponent()}
      </main>
    </div>
  );
};

const App: React.FC = () => (
  <NavigationProvider>
    <NotificationsProvider>
      <ComponentsProvider>
        <DataFetcher>
          <AppContent />
        </DataFetcher>
      </ComponentsProvider>
    </NotificationsProvider>
  </NavigationProvider>
);

const container = document.getElementById("app");
if (!container) {
  throw new Error("Failed to find the root element");
}
const root = createRoot(container);
root.render(<App />);

if ("serviceWorker" in navigator) {
  navigator.serviceWorker.register("/service-worker.js");
}

export default App;
