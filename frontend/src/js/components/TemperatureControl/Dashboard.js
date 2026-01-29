import React from "react";

import { Components } from "../Components/Components";
import { Component } from "../Components/Component";
import { useComponents } from "../ComponentsContext";

import { Card } from "../ui/Card";

export const Dashboard = () => {
  const { components } = useComponents();

  const findComponentByType = (type) => {
    return Object.values(components).find((c) => c.type === type);
  };

  const boiler = findComponentByType("boiler");
  const temperatureSwitch = findComponentByType("homed_temperature_switch");

  return (
    <>
      <div
        className="flex-row"
        style={{
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
      </div>
      <Components typesFilter={["homed_temperature"]} noCard />
    </>
  );
};
