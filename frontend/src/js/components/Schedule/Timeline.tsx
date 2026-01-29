import type React from "react";

import { TimeSlot } from "./TimeSlot";

interface TimeSlotData {
  id: string;
  start: string;
  stop?: string;
  value: number;
  on: boolean;
}

interface TimelineProps {
  day: number;
  data: TimeSlotData[];
  refresh: () => void;
}

export const Timeline: React.FC<TimelineProps> = ({
  day,
  data = [],
  refresh,
}) => {
  if (data.length === 0) {
    return <div>No schedule defined</div>;
  }

  return (
    <div
      style={{
        width: "100%",
        height: "8em",
        overflow: "auto",
        backgroundColor: "#69c0ff",
        borderRadius: "0.3em",
      }}
    >
      <div
        style={{
          height: "100%",
          display: "flex",
          flexDirection: "row",
          justifyContent: "space-around",
        }}
      >
        {data.map((v, i) => (
          <TimeSlot key={i} day={day} refresh={refresh} {...v} />
        ))}
      </div>
    </div>
  );
};
