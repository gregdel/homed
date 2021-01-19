import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";

import { HomedTemperature } from "../HomedComponents/HomedTemperature";
import { Boiler } from "../HomedComponents/Boiler";

import { Row, Col, Card } from "antd";

export const Dashboard = () => {
  const rooms = useSelector((state) => state.stuff.temperatureControl.rooms);
  const boiler = useSelector((state) => state.stuff.temperatureControl.boiler);

  var items = [];
  rooms.forEach((uuid, name) => {
    items.push(<Room key={uuid} name={name} uuid={uuid} />);
  });

  return (
    <>
      <Row style={{ marginBottom: "0.5em" }}>
        <Card title="Boiler" style={{ width: "100%" }}>
          <Boiler uuid={boiler} />
        </Card>
      </Row>
      <Row gutter={[10, 10]} justify="space-around">
        {items}
      </Row>
    </>
  );
};

const Room = ({ uuid }) => (
  <Col xs={24} sm={12} lg={8}>
    <HomedTemperature uuid={uuid} />
  </Col>
);

Room.propTypes = {
  uuid: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
};
