import React, { createContext, useContext, useState, useEffect } from "react";

// Navigation Context and Hook
const NavigationContext = createContext();

const useNavigation = () => {
  // Get path from hash, removing the '#' character
  const getPathFromHash = () => window.location.hash.slice(1) || "/temperature";

  const [currentPath, setCurrentPath] = useState(getPathFromHash());
  const [params, setParams] = useState({});

  useEffect(() => {
    const handleLocationChange = () => {
      const path = getPathFromHash();
      setCurrentPath(path);
      // Extract URL parameters
      const matches = path.match(/\/components\/([^/]+)\/(schedule|graph)/);
      if (matches) {
        setParams({ componentId: matches[1] });
      } else {
        setParams({});
      }
    };

    // Listen to hashchange instead of popstate
    window.addEventListener("hashchange", handleLocationChange);
    // Initial parameter extraction
    handleLocationChange();

    return () => window.removeEventListener("hashchange", handleLocationChange);
  }, []);

  const navigate = (path) => {
    console.warn("navigating to " + path);
    // Update hash instead of using pushState
    window.location.hash = path;
  };

  return { currentPath, params, navigate };
};

// Navigation Provider Component
export const NavigationProvider = ({ children }) => {
  const navigation = useNavigation();
  return (
    <NavigationContext.Provider value={navigation}>
      {children}
    </NavigationContext.Provider>
  );
};

// Custom hook to use navigation in components
export const useNav = () => {
  const context = useContext(NavigationContext);
  if (!context) {
    throw new Error("useNav must be used within a NavigationProvider");
  }
  return context;
};

export const Link = ({
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

  const handleClick = (e) => {
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
