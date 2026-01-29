import { useEffect, useState, useCallback } from "react";
import type { ComponentJSON } from "../types";
import { useComponents } from "./ComponentsContext";

export const WsHandler: React.FC = () => {
  const { replaceComponent } = useComponents();
  const [ws, setWs] = useState<WebSocket | null>(null);

  const stop = useCallback(() => {
    if (!ws) {
      return;
    }
    ws.close();
    setWs(null);
  }, [ws]);

  const connect = useCallback(() => {
    if (ws) {
      return;
    }

    const type = location.protocol === "https:" ? "wss" : "ws";
    const socket = new WebSocket(type + ":" + location.host + "/events");

    socket.onmessage = (event: MessageEvent) => {
      if (!event.data) {
        return;
      }

      const data = JSON.parse(event.data as string) as ComponentJSON;
      if (!data) {
        return;
      }

      replaceComponent(data);
    };

    socket.onerror = () => {
      stop();
    };

    setWs(socket);
  }, [ws, stop, replaceComponent]);

  useEffect(() => {
    const intervalID = setInterval(() => {
      if (!ws) {
        connect();
        return;
      }

      if (ws.readyState === WebSocket.CLOSED) {
        stop();
      }
    }, 10000);

    if (!ws) {
      connect();
    }

    return () => {
      if (ws) {
        stop();
      }
      clearInterval(intervalID);
    };
  }, [ws, connect, stop]);

  return null;
};
