import React, { useEffect } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";

import { componentsFetch } from "../actions/components";

import { WsHandler } from "./Websocket";

export const DataFetcher = ({ children }) => {
  const dispatch = useDispatch();

  const fetchData = () => {
    dispatch(componentsFetch());
  };

  useEffect(() => {
    fetchData();
    // fetch data every time we regain focus
    window.addEventListener("focus", fetchData);
    return () => {
      window.removeEventListener("focus", fetchData);
    };
  }, [dispatch]);

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
