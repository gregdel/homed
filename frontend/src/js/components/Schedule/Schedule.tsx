import type React from "react";
import { useCallback, useEffect, useState } from "react";
import { prettyName } from "../../utils";
import { useNav } from "./../Navigation";

import type { Schedule as ScheduleType } from "../../types";
import { apiGet } from "../../utils/api";
import { Icon } from "../ui/Icon";

import { DailySchedule } from "./DailySchedule";
import { DefaultValue } from "./DefaultValue";
import { Overrides } from "./Overrides";
import { TimeSlotModal } from "./TimeSlotModal";

interface TimeSlotData {
  id: string;
  start: string;
  stop?: string;
  value: number;
  on: boolean;
}

interface ScheduleData {
  schedule_name: string;
  schedule: ScheduleType & {
    days: Record<number, TimeSlotData[]>;
    overrides: OverrideData[];
  };
}

interface OverrideData {
  id: string;
  start: string;
  stop?: string;
  value: number;
  on: boolean;
}

export const Schedule: React.FC = () => {
  const [schedule, setSchedule] = useState<
    ScheduleData["schedule"] | undefined
  >(undefined);
  const [name, setName] = useState<string>("");
  const { params } = useNav();

  const [open, setOpen] = useState<boolean>(false);

  const fetchSchedule = useCallback(async () => {
    try {
      const response = await apiGet<ScheduleData>(
        `/components/${params.componentId}/schedule`,
      );
      setName(response.data.schedule_name);
      setSchedule(response.data.schedule);
    } catch (err) {
      console.error("Error fetching schedule:", err);
    }
  }, [setSchedule, params]);

  useEffect(() => {
    void fetchSchedule();
  }, [fetchSchedule, params]);

  if (schedule === undefined) {
    return null;
  }

  return (
    <>
      <h1>Schedule: {prettyName(name)}</h1>
      <hr className="divider" />
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "space-evenly",
          alignContent: "center",
        }}
      >
        <DefaultValue
          defaultValue={schedule.default_value}
          refresh={() => void fetchSchedule()}
        />
        <span className="divider-vertical" style={{ height: "5rem" }} />
        <TimeSlotModal
          refresh={() => void fetchSchedule()}
          open={open}
          setOpen={setOpen}
        />
        <div
          onClick={() => {
            setOpen(true);
          }}
          style={{ cursor: "pointer", alignSelf: "center" }}
        >
          <div>
            <Icon name="calendarPlus" size={2} />
          </div>
        </div>
      </div>
      <hr className="divider" />
      <Overrides
        refresh={() => void fetchSchedule()}
        overrides={schedule?.overrides || []}
      />
      <hr className="divider" />
      {[1, 2, 3, 4, 5, 6, 0].map((day) => (
        <DailySchedule
          key={day}
          day={day}
          data={schedule?.days?.[day] || []}
          refresh={() => void fetchSchedule()}
        />
      ))}
    </>
  );
};
