import { useEffect, useState, useCallback } from "react";
import { useDispatch } from "react-redux";
import { eventComponentUpdate } from "./actions/homedComponents";

export const WsHandler = () => {
  const dispatch = useDispatch();
  const [ws, setWs] = useState(null);

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

    socket.onmessage = (event) => {
      if (!event.data) {
        return;
      }

      var data = JSON.parse(event.data);
      if (!data) {
        return;
      }

      dispatch(eventComponentUpdate(data));
    };

    socket.onerror = () => {
      stop();
    };

    setWs(socket);
  }, [ws, dispatch, stop]);

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
