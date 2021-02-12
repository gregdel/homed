import React from "react";
import { useDispatch, useSelector } from "react-redux";
import PropTypes from "prop-types";

import { mdiPower } from "@mdi/js";

import { IconToggle } from "./common/IconToggle";

import { componentUpdate } from "../../actions/components";

export const TasmotaSwitch = ({ id }) => {
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
    <IconToggle iconOn={mdiPower} iconOff={mdiPower} toggle={toggle} on={on} />
  );
};

TasmotaSwitch.propTypes = {
  id: PropTypes.string.isRequired,
};
