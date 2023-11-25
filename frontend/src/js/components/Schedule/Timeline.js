import React from "react";
import PropTypes from "prop-types";

import { TimeSlot } from "./TimeSlot";

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
        backgroundColor: "#69c0ff",
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
