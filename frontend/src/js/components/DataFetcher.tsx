import type React from "react";
import { useCallback, useEffect, useRef } from "react";
import type { ReactNode } from "react";

import { useComponents } from "./ComponentsContext";
import { WsHandler } from "./Websocket";
import { resumeRefreshEvent } from "./useResumeRefresh";

interface DataFetcherProps {
  children?: ReactNode;
}

export const DataFetcher: React.FC<DataFetcherProps> = ({ children }) => {
  const { refresh } = useComponents();
  const lastResumeRef = useRef(0);

  const handleResume = useCallback(() => {
    if (document.visibilityState !== "visible") {
      return;
    }

    const now = Date.now();
    if (now - lastResumeRef.current < 500) {
      return;
    }
    lastResumeRef.current = now;

    window.dispatchEvent(new Event(resumeRefreshEvent));
    refresh({ force: true, silent: true });
  }, [refresh]);

  useEffect(() => {
    // Initial fetch
    refresh({ silent: true });

    window.addEventListener("focus", handleResume);
    window.addEventListener("online", handleResume);
    window.addEventListener("pageshow", handleResume);
    document.addEventListener("resume", handleResume);

    // Handle visibility change (works for both tab switching and PWA)
    const handleVisibilityChange = () => {
      if (document.visibilityState === "visible") {
        handleResume();
      }
    };
    document.addEventListener("visibilitychange", handleVisibilityChange);

    return () => {
      window.removeEventListener("focus", handleResume);
      window.removeEventListener("online", handleResume);
      window.removeEventListener("pageshow", handleResume);
      document.removeEventListener("resume", handleResume);
      document.removeEventListener("visibilitychange", handleVisibilityChange);
    };
  }, [refresh, handleResume]);

  return (
    <>
      <WsHandler />
      {children}
    </>
  );
};
