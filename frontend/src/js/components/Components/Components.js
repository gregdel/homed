import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

import { Component } from "./Component";

import { Row, Col } from "antd";

export const Components = ({ typesFilter, noCard }) => {
  const components = useSelector((state) =>
    [...state.components.components]
      .filter(([, component]) => {
        if (!component) {
          return false;
        }

        if (typesFilter.length === 0) {
          return true;
        }

        return typesFilter.includes(component.type);
      })
      .map(([, value]) => ({
        id: value.values.id,
        roomName: value.values.room_name,
        deviceName: value.values.device_name,
        updatedAt: value.values.updated_at,
        friendlyName: value.values.friendly_name,
        type: value.type,
      }))
      .sort((a, b) =>
        a.roomName < b.roomName ? -1 : a.roomName > b.roomName ? 1 : 0
      )
  );

  return (
    <Row gutter={[10, 10]}>
      {components.map(
        ({ id, roomName, deviceName, friendlyName, updatedAt, type }) => (
          <Col key={id} xs={24} sm={12} lg={8}>
            <Component
              id={id}
              type={type}
              noCard={noCard}
              deviceName={deviceName}
              roomName={roomName}
              friendlyName={friendlyName}
              updatedAt={updatedAt}
            />
          </Col>
        )
      )}
    </Row>
  );
};
Components.defaultProps = {
  typesFilter: [],
  noCard: false,
};
Components.propTypes = {
  typesFilter: PropTypes.array.isRequired,
  noCard: PropTypes.bool.isRequired,
};
