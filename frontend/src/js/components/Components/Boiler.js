import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

import Icon from "@mdi/react";
import { mdiFire } from "@mdi/js";

import { Typography } from "antd";

export const Boiler = ({ id }) => {
  const state = useSelector((state) => state.components.components.get(id));
  if (!state) {
    return null;
  }

  const on = state.values.on;

  return (
    <div
      style={{
        display: "flex",
        alignItems: "flex-end",
        flexWrap: "wrap",
        justifyContent: "space-between",
      }}
    >
      <div>
        <Typography.Title level={2}>Boiler</Typography.Title>
      </div>
      <div>
        <div>
          <Icon
            path={mdiFire}
            size={3}
            style={{
              color: on ? "#ff4d4f" : "#00000040",
              transition: "color 0.3s ease-out 0s",
            }}
          />
        </div>
      </div>
    </div>
  );
};

Boiler.propTypes = {
  id: PropTypes.string.isRequired,
};
