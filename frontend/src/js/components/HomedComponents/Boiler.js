import React from "react";
import { useDispatch } from "react-redux";
import PropTypes from "prop-types";

import { Col, Row, Switch } from "antd";

import Icon from "@mdi/react";
import { mdiWaterBoiler } from "@mdi/js";

import { componentUpdate } from "../../actions/homedComponents";

export const Boiler = ({ uuid, on }) => {
  const dispatch = useDispatch();

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
  on: PropTypes.bool.isRequired,
};
