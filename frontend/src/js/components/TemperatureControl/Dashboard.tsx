import type React from "react";

import { Component } from "../Components/Component";
import { Components } from "../Components/Components";
import { useComponents } from "../ComponentsContext";

import type { ComponentJSON } from "../../types";
import { Card } from "../ui/Card";

export const Dashboard: React.FC = () => {
  const { components } = useComponents();

  const findComponentByType = (type: string): ComponentJSON | undefined => {
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
          <Card style={{ marginLeft: "0.5rem", flex: "0 0 auto" }}>
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
