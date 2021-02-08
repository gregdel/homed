import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

export const DeviceStatus = ({ id }) => {
  const online = useSelector(
    (state) => state.components.components.get(id).values.on
  );
  return <>{online ? "Online" : "Offline"}</>;
};

DeviceStatus.propTypes = {
  id: PropTypes.string.isRequired,
};
