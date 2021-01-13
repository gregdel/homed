import React, { useState } from "react";
import PropTypes from "prop-types";
import { useSelector, useDispatch } from "react-redux";

import { componentUpdate } from "../../actions/homedComponents";

import { Typography, Slider } from "antd";
const { Title } = Typography;

export const HomedTemperature = ({ uuid }) => {
  const dispatch = useDispatch();
  const data = useSelector((state) => state.stuff.components.get(uuid));
  if (!data) {
    return null;
  }

  const [newTarget, setNewTarget] = useState(data.values.target);

  const onChange = (value) => {
    setNewTarget(value);
  };

  const onAfterChange = (value) => {
    dispatch(componentUpdate(uuid, { target: value }));
  };

  return (
    <>
      <Title level={1}>{data.values.current}°C</Title>
      <Title level={4}>Set to {data.values.target}°C</Title>
      <Title level={4}>Mode {data.values.mode}</Title>
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
