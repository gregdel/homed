import React, { useState } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";
import { useNav } from "./../Navigation";

import Icon from "@mdi/react";
import { mdiCheckboxBlank, mdiCheckboxBlankOutline } from "@mdi/js";

import { Modal, Form, Input, Switch, Typography, Divider } from "antd";

import {
  addScheduleTimeSlot,
  updateScheduleTimeSlot,
} from "../../actions/schedule";

import { daysMap } from "./DailySchedule";

const DaysCheckboxes = ({ checked, setChecked }) => {
  const handleCheck = (i) => {
    let values = [...checked];
    values[i] = !values[i];
    setChecked(values);
  };

  let items = [];
  for (let i = 0; i < 7; i = i + 1) {
    items.push(
      <div
        style={{
          cursor: "pointer",
          display: "flex",
          justifyContent: "space-between",
          marginLeft: "3em",
          marginRight: "3em",
        }}
        onClick={() => handleCheck(i)}
      >
        <Typography.Title style={{ fontWeight: "300" }} level={3}>
          {daysMap[i]}
        </Typography.Title>
        <div style={{ display: "flex", alignItems: "center" }}>
          <div>
            <Icon
              path={checked[i] ? mdiCheckboxBlank : mdiCheckboxBlankOutline}
              size={2}
            />
          </div>
        </div>
      </div>
    );
  }
  items.push(items.shift());

  return (
    <div style={{ display: "flex", flexDirection: "column" }}>
      <Divider />
      {items}
    </div>
  );
};
DaysCheckboxes.propTypes = {
  checked: PropTypes.array.isRequired,
  setChecked: PropTypes.func.isRequired,
};

export const TimeSlotModal = ({
  start: defaultStart,
  stop: defaultStop,
  target: defaultTarget = 16,
  on: defaultOn = false,
  edit = false,
  open,
  setOpen,
  id,
  day,
}) => {
  const dispatch = useDispatch();
  const { params } = useNav();

  const [start, setStart] = useState(defaultStart);
  const [stop, setStop] = useState(defaultStop);
  const [target, setTarget] = useState(defaultTarget);
  const [on, setOn] = useState(defaultOn);
  const [checked, setChecked] = useState([
    false,
    false,
    false,
    false,
    false,
    false,
    false,
  ]);

  const handleOk = () => {
    setOpen(false);
    const data = {
      start: start,
      stop: stop ? stop : null,
      value: new Number(target),
      on: on,
    };

    if (edit) {
      dispatch(updateScheduleTimeSlot(params.componentId, day, data, id));
    } else {
      if (day === undefined) {
        checked.map((v, day) => {
          if (v === false) {
            return;
          }
          dispatch(addScheduleTimeSlot(params.componentId, day, data));
        });
      } else {
        dispatch(addScheduleTimeSlot(params.componentId, day, data));
      }
    }
  };

  const handleCancel = () => {
    setOpen(false);
  };

  const title = edit ? "Edit timeslot" : "Add in schedule";

  const dayStr = day ? day : "all";
  const formName = id
    ? `ts-${params.componentId}-${dayStr}-modal-${id}`
    : `ts-${params.componentId}-${dayStr}-modal`;

  return (
    <Modal title={title} open={open} onOk={handleOk} onCancel={handleCancel}>
      <Form
        labelCol={{ span: 5 }}
        name={formName}
        onFinish={handleOk}
        onFinishFailed={handleOk}
        initialValues={{
          start: defaultStart,
          stop: defaultStop,
          target: defaultTarget,
        }}
      >
        <Form.Item
          label="From"
          name="start"
          onChange={(e) => setStart(e.target.value)}
        >
          <Input type="time" />
        </Form.Item>

        <Form.Item
          label="To"
          name="stop"
          onChange={(e) => setStop(e.target.value)}
        >
          <Input type="time" />
        </Form.Item>
        <Form.Item
          label="Target"
          name="target"
          onChange={(e) => setTarget(e.target.value)}
        >
          <Input type="number" step=".5" />
        </Form.Item>
        <Form.Item label="Opportunistic" name="on">
          <Switch
            onChange={() => {
              setOn(!on);
            }}
            checked={on}
          />
        </Form.Item>
        {day === undefined && (
          <DaysCheckboxes checked={checked} setChecked={setChecked} />
        )}
      </Form>
    </Modal>
  );
};
TimeSlotModal.propTypes = {
  open: PropTypes.bool.isRequired,
  setOpen: PropTypes.func.isRequired,
  day: PropTypes.number,
  id: PropTypes.string,
  start: PropTypes.any,
  stop: PropTypes.any,
  target: PropTypes.number,
  on: PropTypes.bool,
  edit: PropTypes.bool,
};
