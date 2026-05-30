import type React from "react";
import { createContext, useContext, useEffect, useRef, useState } from "react";
import type { CSSProperties, MouseEvent, ReactNode } from "react";
import type {
  NavigationContext as NavContextType,
  NavigationParams,
} from "../types";

// Navigation Context and Hook
const NavigationContext = createContext<NavContextType | undefined>(undefined);

// Get path from hash, removing the '#' character
const getPathFromHash = (): string =>
  window.location.hash.slice(1) || "/temperature";

// Extract parameters from path
const getParamsFromPath = (path: string): NavigationParams => {
  const matches = path.match(/\/components\/([^/]+)\/(schedule|graph)/);
  return matches?.[1] ? { componentId: matches[1] } : {};
};

const useNavigation = (): NavContextType => {
  const initialPath = getPathFromHash();
  const [currentPath, setCurrentPath] = useState<string>(initialPath);
  const [params, setParams] = useState<NavigationParams>(
    getParamsFromPath(initialPath),
  );
  const currentPathRef = useRef(initialPath);
  const previousPathRef = useRef<string | null>(null);

  useEffect(() => {
    const handleLocationChange = () => {
      const path = getPathFromHash();
      if (path !== currentPathRef.current) {
        previousPathRef.current = currentPathRef.current;
        currentPathRef.current = path;
      }
      setCurrentPath(path);
      setParams(getParamsFromPath(path));
    };

    window.addEventListener("hashchange", handleLocationChange);
    return () => window.removeEventListener("hashchange", handleLocationChange);
  }, []);

  const navigate = (path: string) => {
    if (path === currentPathRef.current) {
      return;
    }
    window.location.hash = path;
  };

  const goBack = (fallbackPath = "/temperature") => {
    if (previousPathRef.current) {
      window.history.back();
      return;
    }

    navigate(fallbackPath);
  };

  return { currentPath, params, navigate, goBack };
};

// Navigation Provider Component
interface NavigationProviderProps {
  children: ReactNode;
}

export const NavigationProvider: React.FC<NavigationProviderProps> = ({
  children,
}) => {
  const navigation = useNavigation();
  return (
    <NavigationContext.Provider value={navigation}>
      {children}
    </NavigationContext.Provider>
  );
};

// Custom hook to use navigation in components
export const useNav = (): NavContextType => {
  const context = useContext(NavigationContext);
  if (!context) {
    throw new Error("useNav must be used within a NavigationProvider");
  }
  return context;
};

interface LinkProps {
  to: string;
  children: ReactNode;
  className?: string;
  activeClassName?: string;
  style?: CSSProperties;
  activeStyle?: CSSProperties;
  onClick?: ((e: MouseEvent<HTMLAnchorElement>) => void) | null;
}

export const Link: React.FC<LinkProps> = ({
  to,
  children,
  className = "",
  activeClassName = "",
  style = {},
  activeStyle = {},
  onClick = null,
}) => {
  const { currentPath, navigate } = useNav();
  const isActive = currentPath === to;

  const handleClick = (e: MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault();
    if (onClick) onClick(e);
    navigate(to);
  };

  return (
    <a
      href={`#${to}`}
      onClick={handleClick}
      className={`${className} ${isActive ? activeClassName : ""}`}
      style={{
        ...style,
        ...(isActive ? activeStyle : {}),
        cursor: "pointer",
        textDecoration: "none",
      }}
    >
      {children}
    </a>
  );
};
