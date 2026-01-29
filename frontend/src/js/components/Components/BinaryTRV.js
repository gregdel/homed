import React from "react";
import PropTypes from "prop-types";

import { Icon } from "../ui/Icon";
import { getSwitch } from "./common/switch.js";

export const BinaryTRV = ({ id }) => {
  const { on } = getSwitch(id);

  return (
    <>
      <div
        style={{
          display: "flex",
          justifyContent: "center",
        }}
      >
        <Icon
          name="radiator"
          size={3}
          style={{
            color: on ? "#ff00005e" : "#00000040",
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
        <span>TRV is {on ? "open" : "closed"}</span>
      </div>
    </>
  );
};

BinaryTRV.propTypes = {
  id: PropTypes.string.isRequired,
};
