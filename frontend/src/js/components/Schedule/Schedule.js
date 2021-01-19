import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
import PropTypes from "prop-types";

import { fetchSchedule } from "../../actions/schedule";

import { Typography, Divider } from "antd";
const { Title } = Typography;

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
      <Timeline data={data} />
      <Divider />
    </>
  );
};
DailySchedule.propTypes = {
  day: PropTypes.number.isRequired,
  data: PropTypes.array.isRequired,
};

export const Timeline = ({ data = [] }) => {
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
        <TimeSlot key={i} {...v} />
      ))}
    </div>
  );
};
Timeline.propTypes = {
  data: PropTypes.array.isRequired,
};

const prettyNumber = (number) => ("0" + number).slice(-2);

export const TimeSlot = ({ start, stop, value }) => {
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
        <div>
          <span style={{ fontSize: "1.5em" }}>{value}°C</span>
        </div>
        <div>
          <span>
            {start.hour}:{prettyNumber(start.minute)}
          </span>
          {stop && (
            <span>
              {" "}
              - {stop.hour}:{prettyNumber(stop.minute)}
            </span>
          )}
        </div>
      </div>
    </div>
  );
};
TimeSlot.propTypes = {
  start: PropTypes.object.isRequired,
  stop: PropTypes.object,
  value: PropTypes.number.isRequired,
};
