import React, { useEffect, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { useParams } from "react-router-dom";
import { prettyName } from "../../utils";

import { Typography, Divider } from "antd";
const { Title } = Typography;

import Icon from "@mdi/react";
import { mdiCalendarPlus } from "@mdi/js";

import { fetchSchedule } from "../../actions/schedule";

import { DefaultValue } from "./DefaultValue";
import { Overrides } from "./Overrides";
import { DailySchedule } from "./DailySchedule";
import { TimeSlotModal } from "./TimeSlotModal";

export const Schedule = () => {
  const dispatch = useDispatch();
  const { componentId: id } = useParams();

  const [open, setOpen] = useState(false);

  useEffect(() => {
    dispatch(fetchSchedule(id));
  }, [dispatch]);

  const data = useSelector((state) => state.schedules.schedules.get(id));
  if (data === undefined) {
    return null;
  }

  let items = [];
  for (let i = 0; i < 7; i = i + 1) {
    items.push(<DailySchedule key={i} day={i} data={data.schedule.days[i]} />);
  }
  items.push(items.shift());

  return (
    <>
      <Title>Schedule: {prettyName(data.schedule_name)}</Title>
      <Divider />
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "space-evenly",
          alignContent: "center",
        }}
      >
        <DefaultValue defaultValue={data.schedule.default_value} />
        <Divider type="vertical" style={{ height: "5rem" }} />
        <TimeSlotModal open={open} setOpen={setOpen} />
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
      <Overrides />
      <Divider />
      {items}
    </>
  );
};
