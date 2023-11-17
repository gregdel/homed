import React from "react";
import { useSelector } from "react-redux";

import PropTypes from "prop-types";

import Icon from "@mdi/react";
import { mdiRadiator } from "@mdi/js";

export const BinaryTRV = ({ id }) => {
  const data = useSelector((state) => state.components.components.get(id));
  if (!data) {
    return null;
  }

  const on = data.values.on;

  return (
    <>
      <div
        style={{
          display: "flex",
          justifyContent: "center",
        }}
      >
        <Icon
          path={mdiRadiator}
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
