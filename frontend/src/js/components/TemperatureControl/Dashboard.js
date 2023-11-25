import React from "react";
import { useSelector } from "react-redux";

import { Components } from "../Components/Components";
import { Component } from "../Components/Component";

import { Card, Row } from "antd";

export const Dashboard = () => {
  const { boiler, temperatureSwitch } = useSelector((state) => ({
    boiler: state.components.temperatureControl.boiler,
    temperatureSwitch: state.components.temperatureControl.temperatureSwitch,
  }));

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
            <Component type="boiler" id={boiler} noCard />
          </Card>
        )}
        {temperatureSwitch !== undefined && (
          <Card style={{ marginLeft: "1rem", flex: "0 0 auto" }}>
            <Component
              type="homed_temperature_switch"
              id={temperatureSwitch}
              noCard
            />
          </Card>
        )}
      </Row>
      <Components typesFilter={["homed_temperature"]} noCard />
    </>
  );
};
