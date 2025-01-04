import React from "react";
import PropTypes from "prop-types";
import { useComponents } from "../ComponentsContext";

export const WifiSignal = ({ id }) => {
  const { getComponentById } = useComponents();
  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { value } = component.values;
  if (!value) {
    return null;
  }

  return <>{value}dB</>;
};

WifiSignal.propTypes = {
  id: PropTypes.string.isRequired,
};
