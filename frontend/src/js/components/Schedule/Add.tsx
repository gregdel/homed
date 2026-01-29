import React, { useState } from "react";

import { Icon } from "../ui/Icon";

import { TimeSlotModal } from "./TimeSlotModal";

interface AddProps {
  day: number;
  refresh: () => void;
}

export const Add: React.FC<AddProps> = ({ day, refresh }) => {
  const [open, setOpen] = useState<boolean>(false);

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
