import React, { useState } from "react";
import { useNav } from "./../Navigation";
import PropTypes from "prop-types";

import Icon from "@mdi/react";
import { mdiPencilOutline, mdiTrashCanOutline, mdiLeaf } from "@mdi/js";

import { TimeSlotModal } from "./TimeSlotModal";

export const TimeSlot = ({ start, stop, value, on, id, day, refresh }) => {
  const { params } = useNav();
  const [open, setOpen] = useState(false);

  const handleDelete = async () => {
    try {
      const response = await fetch(
        `/components/${params.componentId}/schedule/daily/${day}/${id}`,
        {
          method: "DELETE",
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
    } catch (error) {
      console.error("Error posting data:", error);
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
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
        alignItems: "center",
        marginLeft: "0.2em",
        marginRight: "0.2em",
      }}
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
      <div
        style={{
          width: "100%",
          display: "flex",
          flexDirection: "row",
          justifyContent: "space-between",
        }}
      >
        <div style={{ cursor: "pointer" }} onClick={handleEdit}>
          <Icon path={mdiPencilOutline} size={1} />
        </div>
        <div style={{ cursor: "pointer" }} onClick={handleDelete}>
          <Icon path={mdiTrashCanOutline} size={1} />
        </div>
      </div>
      <div style={{ fontSize: "2.3em" }}>
        {value}°C
        {on && <Icon style={{ marginLeft: "0.2em" }} path={mdiLeaf} size={1} />}
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
