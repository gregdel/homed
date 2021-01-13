import React from "react";
import PropTypes from "prop-types";

export const DeviceStatus = ({ online }) => {
  return <>{online ? "Online" : "Offline"}</>;
};

DeviceStatus.propTypes = {
  online: PropTypes.bool.isRequired,
};
