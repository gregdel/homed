import type React from "react";
import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useRef,
  useState,
} from "react";
import type { ReactNode } from "react";
import type { ComponentCommand, ComponentJSON } from "../types";
import { apiGet, apiPut } from "../utils/api";
import { useNotifications } from "./NotificationsContext";

interface ComponentsContextType {
  components: Record<string, ComponentJSON>;
  loading: boolean;
  error: string | null;
  lastUpdated: Date | null;
  updateComponent: (id: string, data: ComponentCommand) => Promise<void>;
  refresh: () => void;
  getComponentById: (id: string) => ComponentJSON | undefined;
  replaceComponent: (component: ComponentJSON) => void;
}

// Create context
const ComponentsContext = createContext<ComponentsContextType | undefined>(
  undefined,
);

// Custom hook for using the components context
export const useComponents = (): ComponentsContextType => {
  const context = useContext(ComponentsContext);
  if (!context) {
    throw new Error("useComponents must be used within a ComponentsProvider");
  }
  return context;
};

interface ComponentsProviderProps {
  children: ReactNode;
}

export const ComponentsProvider: React.FC<ComponentsProviderProps> = ({
  children,
}) => {
  const { addNotificationOk, addNotificationError } = useNotifications();
  const [components, setComponents] = useState<Record<string, ComponentJSON>>(
    {},
  );
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);
  const loadingRef = useRef<boolean>(false);

  const fetchComponents = useCallback(async () => {
    if (loadingRef.current) {
      return;
    }

    try {
      loadingRef.current = true;
      setLoading(true);
      setError(null);

      const data = await apiGet("/components");
      const c: Record<string, ComponentJSON> = {};
      for (const component of data.data as ComponentJSON[]) {
        c[component.values.id] = component;
      }

      setComponents(c);
      setLastUpdated(new Date());
      addNotificationOk("Components updated");
    } catch (err) {
      const errorMessage = err instanceof Error ? err.toString() : String(err);
      setError(errorMessage);
      console.error("Error fetching components:", err);
      addNotificationError("Error while fetching components");
    } finally {
      loadingRef.current = false;
      setLoading(false);
    }
  }, []);

  // Don't refresh more than one per 100ms - use useMemo to keep it stable
  const refresh = useMemo(() => {
    const lastCall = { current: 0 };
    return () => {
      const now = Date.now();
      if (now - lastCall.current > 100) {
        lastCall.current = now;
        void fetchComponents();
      }
    };
  }, [fetchComponents]);

  const updateComponent = async (id: string, data: ComponentCommand) => {
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
    (id: string): ComponentJSON | undefined => {
      return components[id];
    },
    [components],
  );

  const replaceComponent = useCallback((component: ComponentJSON) => {
    setComponents((prevComponents) => ({
      ...prevComponents,
      [component.values.id]: component,
    }));
  }, []);

  const value: ComponentsContextType = {
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
