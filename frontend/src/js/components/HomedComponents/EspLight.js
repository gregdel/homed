import React from "react";
import { useDispatch, useSelector } from "react-redux";

import PropTypes from "prop-types";

import { Col, Row, Switch } from "antd";

import Icon from "@mdi/react";
import { mdiLightbulb } from "@mdi/js";

import { componentUpdate } from "../../actions/homedComponents";

export const ESPLight = ({ uuid }) => {
  const dispatch = useDispatch();
  const data = useSelector((state) => state.stuff.components.get(uuid));
  if (!data) {
    return null;
  }

  const on = data.values.on;

  const toggle = () => {
    const data = { state: on ? "OFF" : "ON" };
    dispatch(componentUpdate(uuid, data));
  };

  return (
    <Row gutter={2}>
      <Col flex={1}>
        <Icon path={mdiLightbulb} size={1} />
      </Col>
      <Col flex={1}>
        <Switch checked={on} onChange={toggle} />
      </Col>
    </Row>
  );
};

ESPLight.propTypes = {
  uuid: PropTypes.string.isRequired,
};
