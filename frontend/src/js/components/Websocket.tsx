import { useCallback, useEffect, useRef } from "react";
import type { ComponentJSON } from "../types";
import { useComponents } from "./ComponentsContext";
import { resumeRefreshEvent } from "./useResumeRefresh";

const reconnectDelay = (attempt: number): number =>
  Math.min(1000 * 2 ** attempt, 30000);

export const WsHandler: React.FC = () => {
  const { refresh, replaceComponent } = useComponents();
  const socketRef = useRef<WebSocket | null>(null);
  const connectRef = useRef<(() => void) | null>(null);
  const reconnectTimerRef = useRef<number | null>(null);
  const reconnectAttemptRef = useRef(0);

  const clearReconnectTimer = useCallback(() => {
    if (reconnectTimerRef.current !== null) {
      window.clearTimeout(reconnectTimerRef.current);
      reconnectTimerRef.current = null;
    }
  }, []);

  const closeSocket = useCallback(() => {
    const socket = socketRef.current;
    socketRef.current = null;

    if (!socket) {
      return;
    }

    socket.onopen = null;
    socket.onmessage = null;
    socket.onerror = null;
    socket.onclose = null;

    if (
      socket.readyState === WebSocket.CONNECTING ||
      socket.readyState === WebSocket.OPEN
    ) {
      socket.close();
    }
  }, []);

  const scheduleReconnect = useCallback(() => {
    if (
      reconnectTimerRef.current !== null ||
      document.visibilityState !== "visible"
    ) {
      return;
    }

    const delay = reconnectDelay(reconnectAttemptRef.current);
    reconnectAttemptRef.current += 1;

    reconnectTimerRef.current = window.setTimeout(() => {
      reconnectTimerRef.current = null;
      connectRef.current?.();
    }, delay);
  }, []);

  const connect = useCallback(() => {
    const currentSocket = socketRef.current;
    if (
      document.visibilityState !== "visible" ||
      currentSocket?.readyState === WebSocket.CONNECTING ||
      currentSocket?.readyState === WebSocket.OPEN
    ) {
      return;
    }

    clearReconnectTimer();

    const type = location.protocol === "https:" ? "wss" : "ws";
    const socket = new WebSocket(`${type}:${location.host}/events`);
    socketRef.current = socket;

    socket.onopen = () => {
      reconnectAttemptRef.current = 0;
      refresh({ silent: true });
    };

    socket.onmessage = (event: MessageEvent) => {
      if (!event.data) {
        return;
      }

      try {
        const data = JSON.parse(event.data as string) as ComponentJSON;
        if (data) {
          replaceComponent(data);
        }
      } catch (err) {
        console.error("Error parsing websocket message:", err);
      }
    };

    socket.onerror = () => {
      socket.close();
    };

    socket.onclose = () => {
      if (socketRef.current === socket) {
        socketRef.current = null;
        scheduleReconnect();
      }
    };
  }, [clearReconnectTimer, refresh, replaceComponent, scheduleReconnect]);

  connectRef.current = connect;

  useEffect(() => {
    const handleHidden = () => {
      clearReconnectTimer();
      closeSocket();
    };

    const handleVisible = () => {
      if (document.visibilityState === "visible") {
        reconnectAttemptRef.current = 0;
        closeSocket();
        connect();
      } else {
        handleHidden();
      }
    };

    const handleResume = () => {
      reconnectAttemptRef.current = 0;
      closeSocket();
      connect();
    };

    connect();
    window.addEventListener("pagehide", handleHidden);
    window.addEventListener(resumeRefreshEvent, handleResume);
    document.addEventListener("freeze", handleHidden);
    document.addEventListener("visibilitychange", handleVisible);

    return () => {
      clearReconnectTimer();
      closeSocket();
      window.removeEventListener("pagehide", handleHidden);
      window.removeEventListener(resumeRefreshEvent, handleResume);
      document.removeEventListener("freeze", handleHidden);
      document.removeEventListener("visibilitychange", handleVisible);
    };
  }, [clearReconnectTimer, closeSocket, connect]);

  return null;
};
