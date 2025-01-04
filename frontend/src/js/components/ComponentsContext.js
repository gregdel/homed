import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
} from "react";
import PropTypes from "prop-types";

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
  const [components, setComponents] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [lastUpdated, setLastUpdated] = useState(null);

  const fetchComponents = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch("/components");
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      if (data.status === "success" && Array.isArray(data.data)) {
        let c = {};
        for (const component of data.data) {
          c[component.values.id] = component;
        }

        setComponents(c);
        setLastUpdated(new Date());
      } else {
        throw new Error("Invalid data format received");
      }
    } catch (err) {
      setError(err.message);
      console.error("Error fetching components:", err);
    } finally {
      setLoading(false);
    }
  }, [setLoading, setError, setComponents, setLastUpdated]);

  const updateComponent = async (id, data) => {
    try {
      const body = typeof data === "string" ? data : JSON.stringify(data);
      const response = await fetch(`/components/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: body,
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
    } catch (error) {
      console.error("Error posting data:", error);
    }
  };

  // Initial fetch on mount
  useEffect(() => {
    fetchComponents();
  }, [fetchComponents]);

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
    refresh: fetchComponents,
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
