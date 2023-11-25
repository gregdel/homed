import React from "react";
import { useDispatch } from "react-redux";
import { useParams } from "react-router-dom";
import PropTypes from "prop-types";

import { deleteSchedule } from "../../actions/schedule";

import Icon from "@mdi/react";
import { mdiTrashCanOutline, mdiLeaf } from "@mdi/js";

// Remove the seconds from the displayed time
const formatTime = (time) => time.slice(0, -3);

export const TimeSlot = ({ start, stop, value, on, id, day }) => {
  const dispatch = useDispatch();
  const { componentId } = useParams();

  const handleDelete = () => {
    dispatch(deleteSchedule(componentId, day, id));
  };

  const color = on ? "#b7eb8f" : "#ffd666";

  return (
    <div
      style={{
        minWidth: "8em",
        backgroundColor: color,
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
      <div style={{ fontSize: "2.3em" }}>
        {value}°C
        {on && <Icon style={{ marginLeft: "0.2em" }} path={mdiLeaf} size={1} />}
      </div>
      <div>
        <span>{formatTime(start)}</span>
        {stop && <span> - {formatTime(stop)}</span>}
      </div>
    </div>
  );
};
TimeSlot.propTypes = {
  id: PropTypes.string.isRequired,
  on: PropTypes.bool.isRequired,
  start: PropTypes.string.isRequired,
  stop: PropTypes.string,
  value: PropTypes.number.isRequired,
  day: PropTypes.number.isRequired,
};
