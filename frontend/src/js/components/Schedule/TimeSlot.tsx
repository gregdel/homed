import React, { useState } from "react";
import { useNav } from "./../Navigation";

import { Icon } from "../ui/Icon";
import { apiDelete } from "../../utils/api";

import { TimeSlotModal } from "./TimeSlotModal";

interface TimeSlotProps {
  start: string;
  stop?: string;
  value: number;
  on: boolean;
  id: string;
  day: number;
  refresh: () => void;
}

export const TimeSlot: React.FC<TimeSlotProps> = ({
  start,
  stop,
  value,
  on,
  id,
  day,
  refresh,
}) => {
  const { params } = useNav();
  const [open, setOpen] = useState<boolean>(false);

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
  const formatTime = (time: string): string => time.slice(0, -3);

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
        {...(stop !== undefined && { stop })}
        on={on}
        id={id}
        target={value}
        refresh={refresh}
        edit
      />
      <div style={{ width: "100%" }} className="flex justify-between">
        <div
          className="cursor-pointer"
          onClick={() => {
            handleEdit();
          }}
        >
          <Icon name="pencilOutline" size={1} />
        </div>
        <div
          className="cursor-pointer"
          onClick={() => {
            void handleDelete();
          }}
        >
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
