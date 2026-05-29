import type React from "react";
import type { RollerShutterValues } from "../../types";
import { useComponents } from "../ComponentsContext";

import { Icon } from "../ui/Icon";
import { IconRollerShutter } from "./common/IconRollerShutter";

interface RollerShutterProps {
  id: string;
}

export const RollerShutter: React.FC<RollerShutterProps> = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const component = getComponentById(id);
  const percentOpen =
    component === undefined
      ? 0
      : (component.values as RollerShutterValues).value;
  const roundedPercentOpen = Math.round(
    Math.max(0, Math.min(100, percentOpen)),
  );

  if (component === undefined) {
    return null;
  }

  const handleClick = (action: string) => {
    void updateComponent(id, action);
  };

  const statusLabel = () => {
    if (roundedPercentOpen === 0) {
      return "Closed";
    }

    if (roundedPercentOpen === 100) {
      return "Open";
    }

    return `Open at ${roundedPercentOpen}%`;
  };

  return (
    <>
      <div className="flex justify-center">
        <IconRollerShutter percentOpen={percentOpen} />

        <div className="flex flex-col justify-center">
          <Icon
            name="arrowUpBoldCircleOutline"
            onClick={() => handleClick("open")}
            className="cursor-pointer"
            size={2}
          />
          <Icon
            name="stopCircleOutline"
            onClick={() => handleClick("stop")}
            className="cursor-pointer"
            size={2}
          />
          <Icon
            name="arrowDownBoldCircleOutline"
            onClick={() => handleClick("close")}
            className="cursor-pointer"
            size={2}
          />
        </div>
      </div>
      <div className="flex justify-center mt-sm text-secondary">
        <span>{statusLabel()}</span>
      </div>
    </>
  );
};
