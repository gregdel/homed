import { useComponents } from "../../ComponentsContext";

export const getSwitch = (id) => {
  const { getComponentById, updateComponent } = useComponents();
  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { on } = component.values;

  const toggle = () => {
    updateComponent(id, on ? "OFF" : "ON");
  };

  return { toggle, on };
};
