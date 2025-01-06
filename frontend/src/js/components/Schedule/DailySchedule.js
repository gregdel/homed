import React from "react";
import PropTypes from "prop-types";

import { Typography, Divider } from "antd";
const { Title } = Typography;

import { Add } from "./Add";
import { Timeline } from "./Timeline";

export const daysMap = {
  0: "Sunday",
  1: "Monday",
  2: "Tuesday",
  3: "Wednesay",
  4: "Thrusday",
  5: "Friday",
  6: "Saturday",
};

export const DailySchedule = ({ day = 0, data = [], refresh }) => {
  return (
    <>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "baseline",
        }}
      >
        <Title level={3}>{daysMap[day]}</Title>
        <Add day={day} refresh={refresh} />
      </div>
      <Timeline day={day} data={data} refresh={refresh} />
      <Divider />
    </>
  );
};
DailySchedule.propTypes = {
  day: PropTypes.number.isRequired,
  data: PropTypes.array.isRequired,
  refresh: PropTypes.func.isRequired,
};
