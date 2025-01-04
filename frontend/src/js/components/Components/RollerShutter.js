import React from "react";
import { useComponents } from "../ComponentsContext";

import PropTypes from "prop-types";

import {
  mdiArrowUpBoldCircleOutline,
  mdiStopCircleOutline,
  mdiArrowDownBoldCircleOutline,
} from "@mdi/js";

import Icon from "@mdi/react";

import { IconRollerShutter } from "./common/IconRollerShutter";

export const RollerShutter = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { value: percentOpen } = component.values;

  const handleClick = (action) => {
    updateComponent(id, action);
  };

  const msg = () => {
    if (percentOpen == 0) {
      return "Closed";
    }

    if (percentOpen == 100) {
      return "Opened";
    }

    return `Opened at ${percentOpen}%`;
  };

  return (
    <>
      <div style={{ display: "flex", justifyContent: "center" }}>
        <IconRollerShutter percentOpen={percentOpen} />

        <div
          style={{
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
          }}
        >
          <Icon
            path={mdiArrowUpBoldCircleOutline}
            onClick={() => handleClick("open")}
            style={{ cursor: "pointer" }}
            size={2}
          />
          <Icon
            path={mdiStopCircleOutline}
            onClick={() => handleClick("stop")}
            style={{ cursor: "pointer" }}
            size={2}
          />
          <Icon
            path={mdiArrowDownBoldCircleOutline}
            onClick={() => handleClick("close")}
            style={{ cursor: "pointer" }}
            size={2}
          />
        </div>
      </div>
      <div
        style={{
          display: "flex",
          justifyContent: "center",
        }}
      >
        <span>{msg()}</span>
      </div>
    </>
  );
};

RollerShutter.propTypes = {
  id: PropTypes.string.isRequired,
};
