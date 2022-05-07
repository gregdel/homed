import React from "react";
import { useSelector } from "react-redux";
import PropTypes from "prop-types";

import { mdiLightningBolt } from "@mdi/js";

import Icon from "@mdi/react";

export const PowerMeter = ({ id }) => {
  const data = useSelector((state) => state.components.components.get(id));
  if (!data) {
    return null;
  }

  return (
    <>
      <div
        style={{
          display: "flex",
          justifyContent: "center",
        }}
      >
        <Icon
          path={mdiLightningBolt}
          size={8}
          style={{
            color: "#ffec3d",
          }}
        />
      </div>
      <div
        style={{
          display: "flex",
          justifyContent: "center",
        }}
      >
        <span>{data.values.value} W</span>
      </div>
    </>
  );
};

PowerMeter.propTypes = {
  id: PropTypes.string.isRequired,
};
