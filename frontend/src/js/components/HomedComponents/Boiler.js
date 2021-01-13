import React from "react";
import { useDispatch, useSelector } from "react-redux";

import PropTypes from "prop-types";

import { Col, Row, Switch } from "antd";

import Icon from "@mdi/react";
import { mdiWaterBoiler } from "@mdi/js";

import { componentUpdate } from "../../actions/homedComponents";

export const Boiler = ({ uuid }) => {
  const dispatch = useDispatch();
  const data = useSelector((state) => state.stuff.components.get(uuid));
  if (!data) {
    return null;
  }

  const on = data.values.on;

  const toggle = () => {
    dispatch(componentUpdate(uuid, on ? "OFF" : "ON"));
  };

  return (
    <Row gutter={2}>
      <Col flex={1}>
        <Icon path={mdiWaterBoiler} size={1} />
      </Col>
      <Col flex={1}>
        <Switch checked={on} onChange={toggle} />
      </Col>
    </Row>
  );
};

Boiler.propTypes = {
  uuid: PropTypes.string.isRequired,
};
