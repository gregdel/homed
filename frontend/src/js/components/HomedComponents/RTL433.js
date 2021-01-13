import React from "react";
import PropTypes from "prop-types";

export const RTL433 = ({ boiler_state }) => {
  return <>{boiler_state ? "Boiler is on" : "Boiler is off"}</>;
};

RTL433.propTypes = {
  boiler_state: PropTypes.bool.isRequired,
};
