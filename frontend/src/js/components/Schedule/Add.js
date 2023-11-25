import React, { useState } from "react";
import PropTypes from "prop-types";

import { Button } from "antd";

import Icon from "@mdi/react";
import { mdiCalendarPlus } from "@mdi/js";

import { TimeSlotModal } from "./TimeSlotModal";

export const Add = ({ day }) => {
  const [open, setOpen] = useState(false);

  const showModal = () => {
    setOpen(true);
  };

  return (
    <div style={{ marginTop: "0.5em" }}>
      <Button
        type="primary"
        size="medium"
        style={{ display: "flex", alignItems: "center" }}
        onClick={showModal}
      >
        <Icon path={mdiCalendarPlus} size={1} />
      </Button>
      <TimeSlotModal day={day} open={open} setOpen={setOpen} />
    </div>
  );
};
Add.propTypes = {
  day: PropTypes.number.isRequired,
};
