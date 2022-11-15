import React from "react";
import { useDispatch, useSelector } from "react-redux";

import PropTypes from "prop-types";

import {
  mdiArrowUpBoldCircleOutline,
  mdiStopCircleOutline,
  mdiArrowDownBoldCircleOutline,
} from "@mdi/js";

import Icon from "@mdi/react";

import { componentUpdate } from "../../actions/components";

export const RollerShutter = ({ id }) => {
  const dispatch = useDispatch();
  const data = useSelector((state) => state.components.components.get(id));
  if (!data) {
    return null;
  }

  const percentOpen = data.values.value;

  const handleClick = (action) => {
    dispatch(componentUpdate(id, action));
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

  // The SVG contains a main frame and 4 panels. Display the main frame in any
  // case and add the panels according to the percentOpen value.
  var svgContent = "M3 4H21V8H19V20H17V8H7V20H5V8H3V4M8";
  if (percentOpen < 100) {
    svgContent += " 9H16V11H8V9M8";
  }
  if (percentOpen <= 66) {
    svgContent += " 12H16V14H8V12M8";
  }
  if (percentOpen <= 33) {
    svgContent += " 15H16V17H8V15M8";
  }
  if (percentOpen == 0) {
    svgContent += " 18H16V20H8V18Z";
  }

  return (
    <>
      <div style={{ display: "flex", justifyContent: "center" }}>
        <svg
          viewBox="0 0 24 24"
          role="presentation"
          style={{
            width: "12rem",
            height: "12rem",
          }}
        >
          <path d={svgContent} style={{ fill: "currentcolor" }}></path>
        </svg>

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
