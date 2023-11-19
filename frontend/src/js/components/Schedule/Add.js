import React, { useState } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";
import { useParams } from "react-router-dom";

import { Modal, Button, Form, Input, Switch } from "antd";

import Icon from "@mdi/react";
import { mdiCalendarPlus } from "@mdi/js";

import { addSchedule } from "../../actions/schedule";

export const Add = ({ day }) => {
  const dispatch = useDispatch();
  const { componentId } = useParams();

  const [show, setShow] = useState(false);
  const [from, setFrom] = useState();
  const [to, setTo] = useState();
  const [target, setTarget] = useState(16);
  const [on, setOn] = useState(false);

  const showModal = () => {
    setShow(true);
  };

  const handleOk = () => {
    setShow(false);
    const data = {
      start: from,
      stop: to ? to : null,
      value: new Number(target),
      on: on,
    };
    dispatch(addSchedule(componentId, day, data));
  };

  const handleCancel = () => {
    setShow(false);
  };

  return (
    <div style={{ marginTop: "0.5em" }}>
      <Button
        type="primary"
        size="medium"
        style={{ display: "flex", alignItems: "center" }}
        onClick={showModal}
      >
        <Icon path={mdiCalendarPlus} size={1} />
      </Button>
      <Modal
        title="Add in schedule"
        open={show}
        onOk={handleOk}
        onCancel={handleCancel}
      >
        <Form
          labelCol={{ span: 5 }}
          name="add-schedule-timeslot"
          onFinish={handleOk}
          onFinishFailed={handleOk}
        >
          <Form.Item
            label="From"
            name="from"
            value={from}
            onChange={(e) => {
              setFrom(e.target.value);
            }}
          >
            <Input type="time" />
          </Form.Item>

          <Form.Item
            label="To"
            name="to"
            value={to}
            onChange={(e) => {
              setTo(e.target.value);
            }}
          >
            <Input type="time" />
          </Form.Item>

          <Form.Item
            label="Target"
            name="target"
            value={target}
            defaultValue={target}
            onChange={(e) => {
              setTarget(e.target.value);
            }}
          >
            <Input type="number" step=".5" />
          </Form.Item>
          <Form.Item label="Opportunistic" name="opportunistic" value={on}>
            <Switch
              checked={on}
              onChange={() => {
                setOn(!on);
              }}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
Add.propTypes = {
  day: PropTypes.number.isRequired,
};
