import React from "react";

import { Components } from "../Components/Components";
import { Component } from "../Components/Component";
import { useComponents } from "../ComponentsContext";

import { Card, Row } from "antd";

export const Dashboard = () => {
  const { components } = useComponents();

  const boiler = Object.values(components).find(
    (c) => c.values.type === "boiler"
  );
  const temperatureSwitch = Object.values(components).find(
    (c) => c.values.type === "homed_temperature_switch"
  );

  return (
    <>
      <Row
        style={{
          marginBottom: "0.5em",
          flexDirection: "row",
          justifyContent: "space-between",
          flexWrap: "nowrap",
          alignItems: "center",
        }}
      >
        {boiler !== undefined && (
          <Card style={{ flex: "1 1 auto" }}>
            <Component type="boiler" id={boiler.values.id} noCard />
          </Card>
        )}
        {temperatureSwitch !== undefined && (
          <Card style={{ marginLeft: "1rem", flex: "0 0 auto" }}>
            <Component
              type="homed_temperature_switch"
              id={temperatureSwitch.values.id}
              noCard
            />
          </Card>
        )}
      </Row>
      <Components typesFilter={["homed_temperature"]} noCard />
    </>
  );
};
