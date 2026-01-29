import type React from "react";
import { useCallback, useEffect } from "react";
import type { ReactNode } from "react";

import { useComponents } from "./ComponentsContext";
import { WsHandler } from "./Websocket";

interface DataFetcherProps {
  children?: ReactNode;
}

export const DataFetcher: React.FC<DataFetcherProps> = ({ children }) => {
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
    const handlePageShow = (event: PageTransitionEvent) => {
      // When navigating to the page from browser cache
      if (event.persisted) {
        refresh();
      }
    };

    window.addEventListener("pageshow", handlePageShow);

    return () => {
      window.removeEventListener("focus", refresh);
      document.removeEventListener("visibilitychange", handleVisibilityChange);
      window.removeEventListener("pageshow", handlePageShow);
    };
  }, [refresh, handleVisibilityChange]);

  return (
    <>
      <WsHandler />
      {children}
    </>
  );
};
