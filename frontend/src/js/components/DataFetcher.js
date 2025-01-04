import React, { useEffect } from "react";
import PropTypes from "prop-types";

import { useComponents } from "./ComponentsContext";
import { WsHandler } from "./Websocket";

export const DataFetcher = ({ children }) => {
  const { refresh } = useComponents();

  useEffect(() => {
    refresh();
    // fetch data every time we regain focus
    window.addEventListener("focus", refresh);
    return () => {
      window.removeEventListener("focus", refresh);
    };
  }, []);

  return (
    <>
      <WsHandler />
      {children}
    </>
  );
};
DataFetcher.propTypes = {
  children: PropTypes.any,
};
