import React, { createContext, useContext, useState, useEffect } from "react";
import PropTypes from "prop-types";

// Navigation Context and Hook
const NavigationContext = createContext();

const useNavigation = () => {
  // Get path from hash, removing the '#' character
  const getPathFromHash = () => window.location.hash.slice(1) || "";

  // Extract parameters from path
  const getParamsFromPath = (path) => {
    const matches = path.match(/\/components\/([^/]+)\/(schedule|graph)/);
    return matches ? { componentId: matches[1] } : {};
  };

  const initialPath = getPathFromHash();
  const [currentPath, setCurrentPath] = useState(initialPath);
  const [params, setParams] = useState(getParamsFromPath(initialPath));

  useEffect(() => {
    const handleLocationChange = () => {
      const path = getPathFromHash();
      setCurrentPath(path);
      setParams(getParamsFromPath(path));
    };

    window.addEventListener("hashchange", handleLocationChange);
    return () => window.removeEventListener("hashchange", handleLocationChange);
  }, []);

  const navigate = (path) => {
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
NavigationProvider.propTypes = {
  children: PropTypes.node,
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

Link.propTypes = {
  children: PropTypes.node,
  to: PropTypes.string.isRequired,
  className: PropTypes.string,
  activeClassName: PropTypes.string,
  style: PropTypes.object,
  activeStyle: PropTypes.object,
  onClick: PropTypes.func,
};
