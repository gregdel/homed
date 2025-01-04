import React from "react";
import { useComponents } from "../ComponentsContext";

import PropTypes from "prop-types";

import Icon from "@mdi/react";
import { mdiRadiator } from "@mdi/js";

export const BinaryTRV = ({ id }) => {
  const { getComponentById } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { on } = component.values;

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
