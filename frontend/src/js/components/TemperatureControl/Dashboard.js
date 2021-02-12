import React from "react";
import { useSelector } from "react-redux";

import { Components } from "../Components/Components";
import { Component } from "../Components/Component";

import { Card, Row } from "antd";

export const Dashboard = () => {
  const boiler = useSelector(
    (state) => state.components.temperatureControl.boiler
  );

  return (
    <>
      <Row style={{ marginBottom: "0.5em" }}>
        <Card style={{ width: "100%" }}>
          <Component type="boiler" id={boiler} noCard />
        </Card>
      </Row>
      <Components typesFilter={["homed_temperature"]} noCard />
    </>
  );
};
