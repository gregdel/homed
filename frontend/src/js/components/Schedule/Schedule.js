import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { useParams } from "react-router-dom";
import PropTypes from "prop-types";
import { prettyName } from "../../utils";

import { fetchSchedule, deleteSchedule } from "../../actions/schedule";

import Icon from "@mdi/react";
import { mdiTrashCanOutline } from "@mdi/js";

import { Typography, Divider } from "antd";
const { Title } = Typography;

import { Add } from "./Add";
import { DefaultValue } from "./DefaultValue";
import { Overrides } from "./Overrides";

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

const DailySchedule = ({ day = 0, data = [] }) => {
  const days = {
    0: "Sunday",
    1: "Monday",
    2: "Tuesday",
    3: "Wednesay",
    4: "Thrusday",
    5: "Friday",
    6: "Saturday",
  };

  return (
    <>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "baseline",
        }}
      >
        <Title level={3}>{days[day]}</Title>
        <Add day={day} />
      </div>
      <Timeline day={day} data={data} />
      <Divider />
    </>
  );
};
DailySchedule.propTypes = {
  day: PropTypes.number.isRequired,
  data: PropTypes.array.isRequired,
};

export const Timeline = ({ day, data = [] }) => {
  if (data.length === 0) {
    return <div>No schedule defined</div>;
  }

  return (
    <div
      style={{
        width: "100%",
        height: "8em",
        overflow: "auto",
        backgroundColor: "#91d5ff",
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
          <TimeSlot key={i} day={day} {...v} />
        ))}
      </div>
    </div>
  );
};
Timeline.propTypes = {
  day: PropTypes.number.isRequired,
  data: PropTypes.array.isRequired,
};

// Remove the seconds from the displayed time
const formatTime = (time) => time.slice(0, -3);

export const TimeSlot = ({ start, stop, value, id, day }) => {
  const dispatch = useDispatch();
  const { componentId } = useParams();

  const handleDelete = () => {
    dispatch(deleteSchedule(componentId, day, id));
  };

  return (
    <div
      style={{
        minWidth: "8em",
        backgroundColor: "#ffd666",
        padding: "0.3em",
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
        alignItems: "center",
        marginLeft: "0.2em",
        marginRight: "0.2em",
      }}
    >
      <div
        style={{ cursor: "pointer", alignSelf: "flex-end" }}
        onClick={handleDelete}
      >
        <Icon path={mdiTrashCanOutline} size={1} />
      </div>
      <div style={{ fontSize: "2.4em" }}>{value}°C</div>
      <div>
        <span>{formatTime(start)}</span>
        {stop && <span> - {formatTime(stop)}</span>}
      </div>
    </div>
  );
};
TimeSlot.propTypes = {
  id: PropTypes.string.isRequired,
  start: PropTypes.string.isRequired,
  stop: PropTypes.string,
  value: PropTypes.number.isRequired,
  day: PropTypes.number.isRequired,
};
