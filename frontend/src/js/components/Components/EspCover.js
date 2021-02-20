import React from "react";
import { useDispatch, useSelector } from "react-redux";

import PropTypes from "prop-types";

import {
  mdiWindowShutter,
  mdiWindowShutterOpen,
  mdiArrowUpBoldCircleOutline,
  mdiStopCircleOutline,
  mdiArrowDownBoldCircleOutline,
} from "@mdi/js";

import Icon from "@mdi/react";

import { componentUpdate } from "../../actions/components";

export const ESPCover = ({ id }) => {
  const dispatch = useDispatch();
  const data = useSelector((state) => state.components.components.get(id));
  if (!data) {
    return null;
  }

  const percentOpen = data.values.value;

  const icon = percentOpen === 100 ? mdiWindowShutterOpen : mdiWindowShutter;

  const handleClick = (action) => {
    dispatch(componentUpdate(id, action));
  };

  return (
    <>
      <div style={{ display: "flex", justifyContent: "center" }}>
        <Icon path={icon} size={8} />
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
        <span>Opened at {percentOpen}%</span>
      </div>
    </>
  );
};

ESPCover.propTypes = {
  id: PropTypes.string.isRequired,
};
