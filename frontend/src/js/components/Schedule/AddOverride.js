import React, { useState } from "react";
import { useDispatch } from "react-redux";
import { useNav } from "./../Navigation";

import { Modal, Button, Form, Input, Switch } from "antd";

import Icon from "@mdi/react";
import { mdiCalendarPlus } from "@mdi/js";

import { addScheduleOverride } from "../../actions/schedule";

export const AddOverride = () => {
  const dispatch = useDispatch();
  const { params } = useNav();

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
    dispatch(addScheduleOverride(params.componentId, data));
  };

  const handleCancel = () => {
    setShow(false);
  };

  return (
    <div style={{ marginTop: "0.5em" }}>
      <Button
        type="primary"
        style={{ margin: 1, display: "flex", alignItems: "center" }}
        onClick={showModal}
      >
        <Icon path={mdiCalendarPlus} size={1} style={{ marginRight: "2px" }} />
      </Button>
      <Modal
        title="Add schedule override"
        open={show}
        onOk={handleOk}
        onCancel={handleCancel}
      >
        <Form
          labelCol={{ span: 5 }}
          name="add-schedule-override"
          onFinish={handleOk}
          onFinishFailed={handleOk}
        >
          <Form.Item
            label="From"
            name="from"
            value={from}
            onChange={(e) => {
              const date = new Date(e.target.value);
              setFrom(date.toISOString());
            }}
          >
            <Input type="datetime-local" />
          </Form.Item>

          <Form.Item
            label="To"
            name="to"
            value={to}
            onChange={(e) => {
              const date = new Date(e.target.value);
              setTo(date.toISOString());
            }}
          >
            <Input type="datetime-local" />
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
