import React, { useState } from "react";
import { useNav } from "./../Navigation";
import PropTypes from "prop-types";

import { Icon } from "../ui/Icon";
import { apiDelete } from "../../utils/api";

import { TimeSlotModal } from "./TimeSlotModal";

export const TimeSlot = ({ start, stop, value, on, id, day, refresh }) => {
  const { params } = useNav();
  const [open, setOpen] = useState(false);

  const handleDelete = async () => {
    try {
      await apiDelete(
        `/components/${params.componentId}/schedule/daily/${day}/${id}`
      );
    } catch (error) {
      console.error("Error deleting time slot:", error);
    } finally {
      refresh();
    }
  };

  const handleEdit = () => {
    setOpen(true);
  };

  // Remove the seconds from the displayed time
  const formatTime = (time) => time.slice(0, -3);

  const color = on ? "#b7eb8f" : "#ffd666";

  if (!id) {
    return null;
  }

  return (
    <div
      style={{
        minWidth: "8em",
        backgroundColor: color,
        padding: "0.3em",
        marginLeft: "0.2em",
        marginRight: "0.2em",
      }}
      className="flex flex-col justify-between items-center"
    >
      <TimeSlotModal
        day={day}
        open={open}
        setOpen={setOpen}
        start={start}
        stop={stop}
        on={on}
        id={id}
        target={value}
        refresh={refresh}
        edit
      />
      <div style={{ width: "100%" }} className="flex justify-between">
        <div className="cursor-pointer" onClick={handleEdit}>
          <Icon name="pencilOutline" size={1} />
        </div>
        <div className="cursor-pointer" onClick={handleDelete}>
          <Icon name="trashCanOutline" size={1} />
        </div>
      </div>
      <div style={{ fontSize: "2.3em" }}>
        {value}°C
        {on && <Icon style={{ marginLeft: "0.2em" }} name="leaf" size={1} />}
      </div>
      <div>
        <span>{formatTime(start)}</span>
        {stop && <span> - {formatTime(stop)}</span>}
      </div>
    </div>
  );
};
TimeSlot.propTypes = {
  id: PropTypes.string.isRequired,
  on: PropTypes.bool.isRequired,
  start: PropTypes.string.isRequired,
  stop: PropTypes.string,
  value: PropTypes.number.isRequired,
  day: PropTypes.number.isRequired,
  refresh: PropTypes.func.isRequired,
};
