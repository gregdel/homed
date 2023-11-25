import React, { useState } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";
import { useParams } from "react-router-dom";

import {
  addScheduleTimeSlot,
  updateScheduleTimeSlot,
} from "../../actions/schedule";

import { Modal, Form, Input, Switch } from "antd";

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
      dispatch(addScheduleTimeSlot(componentId, day, data));
    }
  };

  const handleCancel = () => {
    setOpen(false);
  };

  const title = edit ? "Edit timeslot" : "Add in schedule";

  const formName = id
    ? `ts-${componentId}-${day}-modal-${id}`
    : `ts-${componentId}-${day}-modal`;

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
      </Form>
    </Modal>
  );
};
TimeSlotModal.propTypes = {
  open: PropTypes.bool.isRequired,
  setOpen: PropTypes.func.isRequired,
  day: PropTypes.number.isRequired,
  id: PropTypes.string,
  start: PropTypes.any,
  stop: PropTypes.any,
  target: PropTypes.number,
  on: PropTypes.bool,
  edit: PropTypes.bool,
};
