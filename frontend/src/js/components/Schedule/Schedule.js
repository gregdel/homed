import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import { useParams } from "react-router-dom";
import PropTypes from "prop-types";

import { fetchSchedule, deleteSchedule } from "../../actions/schedule";

import Icon from "@mdi/react";
import { mdiTrashCanOutline } from "@mdi/js";

import { Typography, Divider } from "antd";
const { Title } = Typography;

import { Add } from "./Add";

export const Schedule = () => {
  const dispatch = useDispatch();
  const { componentId: id } = useParams();

  useEffect(() => {
    dispatch(fetchSchedule(id));
  }, [dispatch]);

  const schedule = useSelector((state) => state.schedules.schedules.get(id));
  console.log(schedule);
  if (!schedule || !schedule.days) {
    return null;
  }

  let items = [];
  for (let i = 0; i < 7; i = i + 1) {
    items.push(<DailySchedule key={i} day={i} data={schedule.days[i]} />);
  }
  items.push(items.shift());

  return (
    <>
      <Title>Schedule</Title>
      <Divider />
      {items}
    </>
  );
};

export const DailySchedule = ({ day = 0, data = [] }) => {
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
      <Title level={3}>{days[day]}</Title>
      <Timeline day={day} data={data} />
      <Add day={day} />
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
        height: "6em",
        display: "flex",
        flexDirection: "row",
        alignItems: "flex-end",
        backgroundColor: "#91d5ff",
        borderRadius: "0.3em",
      }}
    >
      {data.map((v, i) => (
        <TimeSlot key={i} day={day} {...v} />
      ))}
    </div>
  );
};
Timeline.propTypes = {
  day: PropTypes.number.isRequired,
  data: PropTypes.array.isRequired,
};

// Remove the seconds from the displayed time
const formatTime = (time) => time.slice(0, -3);

export const TimeSlot = ({ start, stop, value, uuid, day }) => {
  const dispatch = useDispatch();
  const { componentId: id } = useParams();

  const handleDelete = () => {
    dispatch(deleteSchedule(id, day, uuid));
  };

  return (
    <div
      style={{
        height: "5em",
        backgroundColor: "#ffd666",
        width: "8em",
        borderTopLeftRadius: "0.3em",
        borderTopRightRadius: "0.3em",
        marginLeft: "10em",
      }}
    >
      <div
        style={{ display: "flex", flexDirection: "column", margin: "0.3em" }}
      >
        <div style={{ display: "flex", justifyContent: "space-between" }}>
          <div style={{ fontSize: "1.5em" }}>{value}°C</div>
          <div style={{ cursor: "pointer" }} onClick={handleDelete}>
            <Icon path={mdiTrashCanOutline} size={1} />
          </div>
        </div>
        <div>
          <span>{formatTime(start)}</span>
          {stop && <span> - {formatTime(stop)}</span>}
        </div>
      </div>
    </div>
  );
};
TimeSlot.propTypes = {
  uuid: PropTypes.string.isRequired,
  start: PropTypes.string.isRequired,
  stop: PropTypes.string,
  value: PropTypes.number.isRequired,
  day: PropTypes.number.isRequired,
};
