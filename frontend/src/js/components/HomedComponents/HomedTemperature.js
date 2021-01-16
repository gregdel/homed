import React, { useState } from "react";
import PropTypes from "prop-types";
import { useSelector, useDispatch } from "react-redux";

import { componentUpdate } from "../../actions/homedComponents";

import { Slider } from "antd";

export const HomedTemperature = ({ uuid }) => {
  const dispatch = useDispatch();
  const { current, target, mode } = useSelector(
    (state) => state.stuff.components.get(uuid).values
  );

  const [newTarget, setNewTarget] = useState(target);

  const onChange = (value) => {
    setNewTarget(value);
  };

  const onAfterChange = (value) => {
    dispatch(componentUpdate(uuid, { target: value }));
  };

  return (
    <>
      <div style={{ fontSize: "4em" }}>
        <span>{current}°C</span>
      </div>

      <div style={{ fontSize: "1em" }}>
        <span>Heating to {target}°C</span>
      </div>

      <div style={{ fontSize: "1em" }}>
        <span>Mode {mode}</span>
      </div>

      <Slider
        min={5}
        max={35}
        step={0.5}
        value={newTarget}
        onChange={onChange}
        onAfterChange={onAfterChange}
      />
    </>
  );
};

HomedTemperature.propTypes = {
  uuid: PropTypes.string.isRequired,
};
