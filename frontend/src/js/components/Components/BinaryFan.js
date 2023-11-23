import React from "react";
import { useDispatch, useSelector } from "react-redux";

import PropTypes from "prop-types";

import { mdiFan } from "@mdi/js";

import { IconToggle } from "./common/IconToggle";

import { componentUpdate } from "../../actions/components";

export const BinaryFan = ({ id }) => {
  const dispatch = useDispatch();
  const data = useSelector((state) => state.components.components.get(id));
  if (!data) {
    return null;
  }

  const on = data.values.on;

  const toggle = () => {
    dispatch(componentUpdate(id, on ? "OFF" : "ON"));
  };

  return (
    <IconToggle
      iconOn={mdiFan}
      iconOff={mdiFan}
      toggle={toggle}
      on={on}
      rotate={on}
    />
  );
};

BinaryFan.propTypes = {
  id: PropTypes.string.isRequired,
};
