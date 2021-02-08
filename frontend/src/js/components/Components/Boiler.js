import React from "react";
import PropTypes from "prop-types";
import { useDispatch, useSelector } from "react-redux";

import { Col, Row, Switch } from "antd";

import Icon from "@mdi/react";
import { mdiWaterBoiler } from "@mdi/js";

import { componentUpdate } from "../../actions/components";

export const Boiler = ({ id }) => {
  const dispatch = useDispatch();

  const on = useSelector((state) =>
    state.components.components.get(id)
      ? state.components.components.get(id).values.on
      : undefined
  );
  if (on === undefined) {
    return null;
  }

  const toggle = () => {
    dispatch(componentUpdate(id, on ? "OFF" : "ON"));
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
  id: PropTypes.string.isRequired,
};
