import React from "react";
import PropTypes from "prop-types";

import { Icon } from "../../ui/Icon";

export const IconToggle = ({ iconOn, iconOff, toggle, on, rotate = false }) => (
  <div
    style={{
      display: "flex",
      justifyContent: "center",
      flexFlow: "column nowrap",
      height: "30vh",
      WebkitTapHighlightColor: "transparent",
      userSelect: "none",
      outline: "none",
    }}
  >
    <Icon
      name={on ? iconOn : iconOff}
      size={6}
      onClick={toggle}
      className={rotate ? "rotating" : "rotating paused"}
      style={{
        cursor: "pointer",
        color: on ? "#ffec3d" : "#00000040",
        transition: "color 0.3s ease-out 0s",
        alignSelf: "center",
      }}
    />
    <div
      style={{
        alignSelf: "center",
      }}
    >
      {on ? "On" : "Off"}
    </div>
  </div>
);
IconToggle.propTypes = {
  iconOn: PropTypes.string.isRequired,
  iconOff: PropTypes.string.isRequired,
  toggle: PropTypes.func.isRequired,
  on: PropTypes.bool.isRequired,
  rotate: PropTypes.bool,
};
