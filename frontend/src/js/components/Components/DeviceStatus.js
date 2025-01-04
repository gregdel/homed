import React from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

export const DeviceStatus = ({ id }) => {
  const { getComponentById } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { on: online } = component.values;
  return <>{online ? "Online" : "Offline"}</>;
};

DeviceStatus.propTypes = {
  id: PropTypes.string.isRequired,
};
