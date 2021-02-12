import React from "react";
import PropTypes from "prop-types";

import Icon from "@mdi/react";

export const IconToggle = ({ iconOn, iconOff, toggle, on }) => (
  <>
    <div
      style={{
        display: "flex",
        justifyContent: "center",
      }}
    >
      <Icon
        path={on ? iconOn : iconOff}
        size={8}
        onClick={toggle}
        style={{
          cursor: "pointer",
          color: on ? "#ffec3d" : "#00000040",
          transition: "color 0.3s ease-out 0s",
        }}
      />
    </div>
    <div
      style={{
        display: "flex",
        justifyContent: "center",
      }}
    >
      <span>{on ? "On" : "Off"}</span>
    </div>
  </>
);
IconToggle.propTypes = {
  iconOn: PropTypes.string.isRequired,
  iconOff: PropTypes.string.isRequired,
  toggle: PropTypes.func.isRequired,
  on: PropTypes.bool.isRequired,
};
