import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

import { Component } from "./Component";

import { Row, Col } from "antd";

export const Components = ({ typesFilter = [], noCard = false }) => {
  const components = useSelector((state) =>
    [...state.components.components]
      .filter(([, component]) => {
        if (!component) {
          return false;
        }

        if (typesFilter.length === 0) {
          return true;
        }

        if (component.values.hide) {
          return false;
        }

        return typesFilter.includes(component.type);
      })
      .map(([, value]) => ({
        id: value.values.id,
        roomName: value.values.room_name,
        friendlyName: value.values.friendly_name,
        type: value.type,
      }))
      .sort((a, b) => {
        if (a.friendlyName !== "" && b.friendlyName !== "") {
          return a.friendlyName < b.friendlyName
            ? -1
            : a.friendlyName > b.friendlyName
            ? 1
            : 0;
        }

        if (a.roomName !== undefined && b.roomName !== undefined) {
          return a.roomName < b.roomName ? -1 : a.roomName > b.roomName ? 1 : 0;
        }

        return a.id < b.id ? -1 : a.id > b.id ? 1 : 0;
      })
  );

  return (
    <Row gutter={[10, 10]}>
      {components.map(({ id, type }) => (
        <Col key={id} xs={24} sm={12} lg={8}>
          <Component id={id} type={type} noCard={noCard} />
        </Col>
      ))}
    </Row>
  );
};
Components.propTypes = {
  typesFilter: PropTypes.array,
  noCard: PropTypes.bool,
};
