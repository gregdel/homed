import React, { useEffect, useState, useCallback } from "react";
import { prettyName } from "../../utils";
import { useNav } from "./../Navigation";

import { Icon } from "../ui/Icon";
import { apiGet } from "../../utils/api";

import { DefaultValue } from "./DefaultValue";
import { Overrides } from "./Overrides";
import { DailySchedule } from "./DailySchedule";
import { TimeSlotModal } from "./TimeSlotModal";

export const Schedule = () => {
  const [schedule, setSchedule] = useState(undefined);
  const [name, setName] = useState("");
  const { params } = useNav();

  const [open, setOpen] = useState(false);

  const fetchSchedule = useCallback(async () => {
    try {
      const data = await apiGet(`/components/${params.componentId}/schedule`);
      setName(data.data.schedule_name);
      setSchedule(data.data.schedule);
    } catch (err) {
      console.error("Error fetching schedule:", err);
    }
  }, [setSchedule, params]);

  useEffect(() => {
    fetchSchedule();
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
          refresh={fetchSchedule}
        />
        <span className="divider-vertical" style={{ height: "5rem" }} />
        <TimeSlotModal refresh={fetchSchedule} open={open} setOpen={setOpen} />
        <div
          onClick={() => setOpen(true)}
          style={{ cursor: "pointer", alignSelf: "center" }}
        >
          <div>
            <Icon name="calendarPlus" size={2} />
          </div>
        </div>
      </div>
      <hr className="divider" />
      <Overrides refresh={fetchSchedule} overrides={schedule.overrides} />
      <hr className="divider" />
      {[1, 2, 3, 4, 5, 6, 0].map((day) => (
        <DailySchedule
          key={day}
          day={day}
          data={schedule.days[day]}
          refresh={fetchSchedule}
        />
      ))}
    </>
  );
};
