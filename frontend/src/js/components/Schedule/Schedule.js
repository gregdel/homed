import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import PropTypes from "prop-types";
import moment from "moment";

import { fetchSchedule } from "../../actions/schedule";

export const Schedule = () => {
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(fetchSchedule());
  }, [dispatch]);

  const schedule = useSelector((state) => state.temperatureSchedule.schedule);
  if (!schedule || !schedule.days) {
    return null;
  }

  let items = [];
  for (let i = 0; i < 7; i = i + 1) {
    items.push(<DailySchedule key={i} day={i} data={schedule.days[i]} />);
  }

  return (
    <>
      <h1>Schedule</h1>
      {items}
    </>
  );
};

export const DailySchedule = ({ day = 0, data = [] }) => {
  return (
    <>
      <h1>Day: {day}</h1>
      <Timeline data={data} />
    </>
  );
};
DailySchedule.propTypes = {
  day: PropTypes.number.required,
  data: PropTypes.array.required,
};

export const Timeline = ({ data = [] }) => {
  if (data.length === 0) {
    return null;
  }

  console.log(data);

  return (
    <div
      style={{
        width: "100%",
        height: "10em",
        display: "flex",
        flexDirection: "row",
        alignItems: "flex-end",
        backgroundColor: "#002766",
      }}
    >
      {data.map((v, i) => (
        <TimeSlot key={i} {...v} />
      ))}
    </div>
  );
};
Timeline.propTypes = {
  data: PropTypes.array.required,
};

export const TimeSlot = ({ uuid, start, stop, value }) => {
  // <p>{uuid}</p>
  // <p>
  //   {start.hour}:{start.minute}
  // </p>
  // {stop && (
  //   <p>
  //     {stop.hour}:{stop.minute}
  //   </p>
  // )}
  return (
    <div
      style={{
        height: "5em",
        backgroundColor: "#ffc53d",
        width: "5em",
        borderTopLeftRadius: "0.3em",
        borderTopRightRadius: "0.3em",
        marginLeft: "10em",
      }}
    >
      <p>{value}</p>
    </div>
  );
};
TimeSlot.propTypes = {
  uuid: PropTypes.string.required,
  start: PropTypes.object.required,
  stop: PropTypes.object.required,
  value: PropTypes.number.required,
};
