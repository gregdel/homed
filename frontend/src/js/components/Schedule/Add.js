import React, { useState } from "react";
import PropTypes from "prop-types";

import { Icon } from "../ui/Icon";

import { TimeSlotModal } from "./TimeSlotModal";

export const Add = ({ day, refresh }) => {
  const [open, setOpen] = useState(false);

  const showModal = () => {
    setOpen(true);
  };

  return (
    <div style={{ marginTop: "0.5em" }}>
      <button className="btn btn-primary" onClick={showModal}>
        <Icon name="calendarPlus" size={1} />
      </button>
      <TimeSlotModal
        day={day}
        open={open}
        setOpen={setOpen}
        refresh={refresh}
      />
    </div>
  );
};
Add.propTypes = {
  day: PropTypes.number.isRequired,
  refresh: PropTypes.func.isRequired,
};
