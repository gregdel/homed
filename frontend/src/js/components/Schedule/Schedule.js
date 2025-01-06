import React, { useEffect, useState, useCallback } from "react";
import { prettyName } from "../../utils";
import { useNav } from "./../Navigation";

import { Typography, Divider } from "antd";
const { Title } = Typography;

import Icon from "@mdi/react";
import { mdiCalendarPlus } from "@mdi/js";

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
      const response = await fetch(
        `/components/${params.componentId}/schedule`
      );
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      if (data.status === "success") {
        setName(data.data.schedule_name);
        setSchedule(data.data.schedule);
      } else {
        throw new Error("Invalid data format received");
      }
    } catch (err) {
      console.error("Error fetching components:", err);
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
      <Title>Schedule: {prettyName(name)}</Title>
      <Divider />
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
        <Divider type="vertical" style={{ height: "5rem" }} />
        <TimeSlotModal refresh={fetchSchedule} open={open} setOpen={setOpen} />
        <div
          onClick={() => setOpen(true)}
          style={{ cursor: "pointer", alignSelf: "center" }}
        >
          <div>
            <Icon path={mdiCalendarPlus} size={2} />
          </div>
        </div>
      </div>
      <Divider />
      <Overrides refresh={fetchSchedule} overrides={schedule.overrides} />
      <Divider />
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
