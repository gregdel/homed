import React from "react";
import { useComponents } from "../ComponentsContext";

import PropTypes from "prop-types";

import { Icon } from "../ui/Icon";
import { IconRollerShutter } from "./common/IconRollerShutter";

export const RollerShutter = ({ id }) => {
  const { getComponentById, updateComponent } = useComponents();

  const component = getComponentById(id);
  if (component === undefined) {
    return null;
  }

  const { value: percentOpen } = component.values;

  const handleClick = (action) => {
    updateComponent(id, action);
  };

  const msg = () => {
    if (percentOpen == 0) {
      return "Closed";
    }

    if (percentOpen == 100) {
      return "Opened";
    }

    return `Opened at ${percentOpen}%`;
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
      <div className="flex justify-center">
        <span>{msg()}</span>
      </div>
    </>
  );
};

RollerShutter.propTypes = {
  id: PropTypes.string.isRequired,
};
