import React, { useState } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";
import { useNav } from "./../Navigation";

import { Modal, Form, Input } from "antd";

import { setScheduleDefault } from "../../actions/schedule";

export const DefaultValue = ({ defaultValue }) => {
  const dispatch = useDispatch();
  const { params } = useNav();
  const [value, setValue] = useState(defaultValue);

  const handleOk = () => {
    setOpen(false);
    dispatch(setScheduleDefault(params.componentId, value));
  };

  const [open, setOpen] = useState(false);

  return (
    <div style={{ cursor: "pointer" }}>
      <div
        style={{ display: "flex", flexDirection: "column" }}
        onClick={() => setOpen(true)}
      >
        <div style={{ fontSize: "2em" }}>Default:</div>
        <div style={{ fontSize: "3em" }}>{value} °C</div>
      </div>
      <Modal
        title="Set default value"
        open={open}
        onOk={handleOk}
        onCancel={() => setOpen(false)}
      >
        <Form
          labelCol={{ span: 5 }}
          name="schedule-default-value"
          onFinish={handleOk}
          onFinishFailed={handleOk}
          initialValues={{
            defaultValue: defaultValue,
          }}
        >
          <Form.Item
            label="Default value"
            name="defaultValue"
            suffix="°C"
            onChange={(e) => setValue(e.target.value)}
          >
            <Input type="number" step=".5" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
DefaultValue.propTypes = {
  defaultValue: PropTypes.number.isRequired,
};
