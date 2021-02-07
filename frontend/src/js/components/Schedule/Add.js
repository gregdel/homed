import React, { useState } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";
import { useParams } from "react-router-dom";

import { Modal, Button, Form, Input } from "antd";

import { addSchedule } from "../../actions/schedule";

export const Add = ({ day }) => {
  const dispatch = useDispatch();
  const { componentId } = useParams();

  const [show, setShow] = useState(false);
  const [from, setFrom] = useState();
  const [to, setTo] = useState();
  const [target, setTarget] = useState(16);

  const showModal = () => {
    setShow(true);
  };

  const handleOk = () => {
    setShow(false);
    const data = {
      start: from,
      stop: to ? to : null,
      value: new Number(target),
    };
    dispatch(addSchedule(componentId, day, data));
  };

  const handleCancel = () => {
    setShow(false);
  };

  return (
    <>
      <Button type="primary" style={{ margin: 1 }} onClick={showModal}>
        Add
      </Button>
      <Modal
        title="Add in schedule"
        visible={show}
        onOk={handleOk}
        onCancel={handleCancel}
      >
        <Form name="basic" onFinish={handleOk} onFinishFailed={handleOk}>
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
            <Input type="number" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};
Add.propTypes = {
  day: PropTypes.number.isRequired,
};
