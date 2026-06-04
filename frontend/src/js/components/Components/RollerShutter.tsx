import type React from "react";
import type { RollerShutterValues } from "../../types";
import { useComponents } from "../ComponentsContext";

import { Icon } from "../ui/Icon";
import { IconRollerShutter } from "./common/IconRollerShutter";

interface RollerShutterProps {
  id: string;
}

const formatNextEventTime = (scheduledAt: string): string =>
  new Date(scheduledAt).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });

export const RollerShutter: React.FC<RollerShutterProps> = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const component = getComponentById(id);
  const values = component?.values as RollerShutterValues | undefined;
  const percentOpen = values?.value ?? 0;
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
      <div className="flex flex-col items-center mt-sm gap-xs text-secondary">
        <span>{statusLabel()}</span>
        {values?.next_event !== null && values?.next_event !== undefined && (
          <div>
            <span>Next {values.next_event.action}</span>
            <span className="font-semibold">
              {" "}
              {formatNextEventTime(values.next_event.scheduled_at)}
            </span>
          </div>
        )}
      </div>
    </>
  );
};
