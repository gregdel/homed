import React, { useEffect, useCallback } from "react";
import PropTypes from "prop-types";

import { useComponents } from "./ComponentsContext";
import { WsHandler } from "./Websocket";

export const DataFetcher = ({ children }) => {
  const { refresh } = useComponents();

  const handleVisibilityChange = useCallback(() => {
    if (document.visibilityState === "visible") {
      refresh();
    }
  }, [refresh]);

  useEffect(() => {
    // Initial fetch
    refresh();

    // Handle webpage focus
    window.addEventListener("focus", refresh);

    // Handle visibility change (works for both tab switching and PWA)
    document.addEventListener("visibilitychange", handleVisibilityChange);

    // Optional: Handle PWA-specific resume event for iOS
    window.addEventListener("pageshow", (event) => {
      // When navigating to the page from browser cache
      if (event.persisted) {
        refresh();
      }
    });

    return () => {
      window.removeEventListener("focus", refresh);
      document.removeEventListener("visibilitychange", handleVisibilityChange);
      window.removeEventListener("pageshow", refresh);
    };
  }, [refresh, handleVisibilityChange]);

  return (
    <>
      <WsHandler />
      {children}
    </>
  );
};

DataFetcher.propTypes = {
  children: PropTypes.node,
};
