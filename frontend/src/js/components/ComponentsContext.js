import React, {
  createContext,
  useContext,
  useState,
  useCallback,
  useRef,
} from "react";
import PropTypes from "prop-types";

import { useNotifications } from "./NotificationsContext";
import { apiGet, apiPut } from "../utils/api";

// Create context
const ComponentsContext = createContext();

// Custom hook for using the components context
export const useComponents = () => {
  const context = useContext(ComponentsContext);
  if (!context) {
    throw new Error("useComponents must be used within a ComponentsProvider");
  }
  return context;
};

export const ComponentsProvider = ({ children }) => {
  const { addNotificationOk, addNotificationError } = useNotifications();
  const [components, setComponents] = useState({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [lastUpdated, setLastUpdated] = useState(null);

  const fetchComponents = useCallback(async () => {
    if (loading) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const data = await apiGet("/components");
      let c = {};
      for (const component of data.data) {
        c[component.values.id] = component;
      }

      setComponents(c);
      setLastUpdated(new Date());
      addNotificationOk("Components updated");
    } catch (err) {
      setError(err.toString());
      console.error("Error fetching components:", err);
      addNotificationError("Error while fetching components");
    } finally {
      setLoading(false);
    }
  }, [setLoading, setError, setComponents, setLastUpdated]);

  const useThrottle = (callback, delay) => {
    const lastCall = useRef(0);
    return useCallback(() => {
      const now = Date.now();
      if (now - lastCall.current > delay) {
        lastCall.current = now;
        callback();
      }
    }, [callback, delay]);
  };

  // Don't refresh more than one per 100ms
  const refresh = useThrottle(fetchComponents, 100);

  const updateComponent = async (id, data) => {
    try {
      // If data is already a string, send it as-is, otherwise it's an object that needs JSON.stringify
      if (typeof data === "string") {
        await fetch(`/components/${id}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: data,
        }).then((response) => {
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
        });
      } else {
        await apiPut(`/components/${id}`, data);
      }
    } catch (error) {
      console.error("Error updating component:", error);
    }
  };

  const getComponentById = useCallback(
    (id) => {
      return components[id];
    },
    [components]
  );

  const replaceComponent = useCallback((component) => {
    setComponents((prevComponents) => ({
      ...prevComponents,
      [component.values.id]: component,
    }));
  }, []);

  const value = {
    components,
    loading,
    error,
    lastUpdated,
    updateComponent,
    refresh,
    getComponentById,
    replaceComponent,
  };

  return (
    <ComponentsContext.Provider value={value}>
      {children}
    </ComponentsContext.Provider>
  );
};

ComponentsProvider.propTypes = {
  children: PropTypes.node,
};
