import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { useParams } from "react-router-dom";
import { prettyName } from "../../utils";

import { fetchSchedule } from "../../actions/schedule";

import { Typography, Divider } from "antd";
const { Title } = Typography;

import { DefaultValue } from "./DefaultValue";
import { Overrides } from "./Overrides";
import { DailySchedule } from "./DailySchedule";

export const Schedule = () => {
  const dispatch = useDispatch();
  const { componentId: id } = useParams();

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
      <DefaultValue defaultValue={data.schedule.default_value} />
      <Divider />
      <Overrides />
      <Divider />
      {items}
    </>
  );
};
