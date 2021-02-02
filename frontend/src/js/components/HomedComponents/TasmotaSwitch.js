import React from "react";
import { useDispatch, useSelector } from "react-redux";

import PropTypes from "prop-types";

import { Col, Row, Switch } from "antd";

import Icon from "@mdi/react";
import { mdiPower } from "@mdi/js";

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
    <Row gutter={2}>
      <Col flex={1}>
        <Icon path={mdiPower} size={1} />
      </Col>
      <Col flex={1}>
        <Switch checked={on} onChange={toggle} />
      </Col>
    </Row>
  );
};

TasmotaSwitch.propTypes = {
  id: PropTypes.string.isRequired,
};
