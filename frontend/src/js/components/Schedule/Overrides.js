import React from "react";
import { useSelector, useDispatch } from "react-redux";
import { useNav } from "./../Navigation";

import PropTypes from "prop-types";

import { Typography, List } from "antd";
const { Title } = Typography;

import Icon from "@mdi/react";
import { mdiTrashCanOutline, mdiLeaf } from "@mdi/js";

import { AddOverride } from "./AddOverride";

import { deleteScheduleOverride } from "../../actions/schedule";

export const Overrides = () => {
  const { params } = useNav();

  const { overrides } = useSelector(
    (state) => state.schedules.schedules.get(params.componentId).schedule
  );

  return (
    <>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "baseline",
        }}
      >
        <Title level={3}>Overrides</Title>
        <AddOverride />
      </div>
      {overrides.length !== 0 && (
        <List
          bordered
          dataSource={overrides}
          renderItem={(v) => <Override {...v} />}
        />
      )}
      {overrides.length === 0 && <div>No override defined</div>}
    </>
  );
};
Overrides.propTypes = {};

const Override = ({ id, start, stop, value, on }) => {
  const dispatch = useDispatch();
  const { params } = useNav();

  const handleDelete = () => {
    dispatch(deleteScheduleOverride(params.componentId, id));
  };

  const formatDate = (date) => {
    const d = new Date(date);
    return d.toLocaleString();
  };

  return (
    <List.Item style={{ flexWrap: "nowrap" }}>
      <div>
        <Typography.Text>
          {on && <Icon path={mdiLeaf} size={0.5} />}
          <strong>{value}°C</strong> from <strong>{formatDate(start)}</strong>{" "}
          to <strong>{formatDate(stop)}</strong>
        </Typography.Text>
      </div>
      <div
        style={{ cursor: "pointer", alignSelf: "flex-end" }}
        onClick={handleDelete}
      >
        <Icon path={mdiTrashCanOutline} size={1} />
      </div>
    </List.Item>
  );
};
Override.propTypes = {
  id: PropTypes.string.isRequired,
  start: PropTypes.string.isRequired,
  stop: PropTypes.string.isRequired,
  value: PropTypes.number.isRequired,
  on: PropTypes.bool.isRequired,
};
