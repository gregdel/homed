import React, { useState } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";
import { useParams } from "react-router-dom";

import {
  addScheduleTimeSlot,
  updateScheduleTimeSlot,
} from "../../actions/schedule";

import { daysMap } from "./DailySchedule";

import { Modal, Form, Input, Switch, Checkbox } from "antd";

const DaysCheckboxes = ({ checked, setChecked }) => {
  const handleCheck = (i) => {
    let values = [...checked];
    values[i] = !values[i];
    setChecked(values);
  };

  let items = [];
  for (let i = 0; i < 7; i = i + 1) {
    items.push(
      <Form.Item label={daysMap[i]} name={`day-${i}`} key={`day-${i}`}>
        <Checkbox checked={checked[i]} onChange={() => handleCheck(i)} />
      </Form.Item>
    );
  }
  items.push(items.shift());

  return items;
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
  const { componentId } = useParams();

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
      dispatch(updateScheduleTimeSlot(componentId, day, data, id));
    } else {
      if (day === undefined) {
        checked.map((v, day) => {
          if (v === false) {
            return;
          }
          dispatch(addScheduleTimeSlot(componentId, day, data));
        });
      } else {
        dispatch(addScheduleTimeSlot(componentId, day, data));
      }
    }
  };

  const handleCancel = () => {
    setOpen(false);
  };

  const title = edit ? "Edit timeslot" : "Add in schedule";

  const dayStr = day ? day : "all";
  const formName = id
    ? `ts-${componentId}-${dayStr}-modal-${id}`
    : `ts-${componentId}-${dayStr}-modal`;

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
