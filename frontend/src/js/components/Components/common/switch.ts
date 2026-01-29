import { useComponents } from "../../ComponentsContext";
import type { SwitchValues } from "../../../types";

export const getSwitch = (id: string) => {
  const { getComponentById, updateComponent } = useComponents();
  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { on } = component.values as SwitchValues;

  const toggle = () => {
    void updateComponent(id, on ? "OFF" : "ON");
  };

  return { toggle, on };
};
