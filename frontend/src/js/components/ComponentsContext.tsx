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
  hasLoaded: boolean;
  error: string | null;
  lastUpdated: Date | null;
  updateComponent: (id: string, data: ComponentCommand) => Promise<void>;
  refresh: (options?: RefreshOptions) => void;
  getComponentById: (id: string) => ComponentJSON | undefined;
  replaceComponent: (component: ComponentJSON) => void;
}

interface RefreshOptions {
  force?: boolean;
  silent?: boolean;
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
  const [hasLoaded, setHasLoaded] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);
  const loadingRef = useRef<boolean>(false);
  const hasLoadedRef = useRef<boolean>(false);
  const activeRequestRef = useRef<AbortController | null>(null);

  const fetchComponents = useCallback(
    async (options: RefreshOptions = {}) => {
      if (loadingRef.current && !options.force) {
        return;
      }

      if (options.force) {
        activeRequestRef.current?.abort();
      }

      const controller = new AbortController();
      activeRequestRef.current = controller;

      try {
        loadingRef.current = true;
        setLoading(!hasLoadedRef.current);
        setError(null);

        const data = await apiGet<ComponentJSON[]>("/components", {
          cache: "no-store",
          signal: controller.signal,
        });
        const c: Record<string, ComponentJSON> = {};
        for (const component of data.data) {
          c[component.values.id] = component;
        }

        setComponents(c);
        setHasLoaded(true);
        hasLoadedRef.current = true;
        setLastUpdated(new Date());
        if (!options.silent) {
          addNotificationOk("Components updated");
        }
      } catch (err) {
        if (err instanceof DOMException && err.name === "AbortError") {
          return;
        }

        const errorMessage =
          err instanceof Error ? err.toString() : String(err);
        setError(errorMessage);
        console.error("Error fetching components:", err);
        addNotificationError("Error while fetching components");
      } finally {
        if (activeRequestRef.current === controller) {
          activeRequestRef.current = null;
          loadingRef.current = false;
          setLoading(false);
        }
      }
    },
    [addNotificationOk, addNotificationError],
  );

  // Don't refresh more than one per 100ms - use useMemo to keep it stable
  const refresh = useMemo(() => {
    const lastCall = { current: 0 };
    return (options: RefreshOptions = {}) => {
      const now = Date.now();
      if (options.force || now - lastCall.current > 100) {
        lastCall.current = now;
        void fetchComponents(options);
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
    setLastUpdated(new Date());
  }, []);

  const value: ComponentsContextType = {
    components,
    loading,
    hasLoaded,
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
