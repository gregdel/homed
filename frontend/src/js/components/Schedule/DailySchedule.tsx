import React from "react";

import { Add } from "./Add";
import { Timeline } from "./Timeline";

export const daysMap: Record<number, string> = {
  0: "Sunday",
  1: "Monday",
  2: "Tuesday",
  3: "Wednesday",
  4: "Thursday",
  5: "Friday",
  6: "Saturday",
};

interface TimeSlotData {
  id: string;
  start: string;
  stop?: string;
  value: number;
  on: boolean;
}

interface DailyScheduleProps {
  day: number;
  data: TimeSlotData[];
  refresh: () => void;
}

export const DailySchedule: React.FC<DailyScheduleProps> = ({
  day = 0,
  data = [],
  refresh,
}) => {
  return (
    <>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "baseline",
        }}
      >
        <h3>{daysMap[day]}</h3>
        <Add day={day} refresh={refresh} />
      </div>
      <Timeline day={day} data={data} refresh={refresh} />
      <hr className="divider" />
    </>
  );
};
