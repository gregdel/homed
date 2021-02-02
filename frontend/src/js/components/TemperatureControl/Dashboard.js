import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

import { HomedTemperature } from "../HomedComponents/HomedTemperature";
import { Boiler } from "../HomedComponents/Boiler";

import { Row, Col, Card } from "antd";

export const Dashboard = () => {
  const rooms = useSelector(
    (state) => state.components.temperatureControl.rooms
  );
  const boiler = useSelector(
    (state) => state.components.temperatureControl.boiler
  );

  var items = [];
  rooms.forEach((id, name) => {
    items.push(<Room key={id} name={name} id={id} />);
  });

  return (
    <>
      <Row style={{ marginBottom: "0.5em" }}>
        <Card title="Boiler" style={{ width: "100%" }}>
          <Boiler id={boiler} />
        </Card>
      </Row>
      <Row gutter={[10, 10]} justify="space-around">
        {items}
      </Row>
    </>
  );
};

const Room = ({ id }) => (
  <Col xs={24} sm={12} lg={8}>
    <HomedTemperature id={id} />
  </Col>
);

Room.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
};
