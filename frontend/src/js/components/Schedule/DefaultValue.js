import React, { useState } from "react";
import PropTypes from "prop-types";
import { useNav } from "./../Navigation";

import { Modal, Form, Input } from "antd";

export const DefaultValue = ({ defaultValue, refresh }) => {
  const { params } = useNav();
  const [value, setValue] = useState(defaultValue);

  const setDefault = async () => {
    try {
      const response = await fetch(
        `/components/${params.componentId}/schedule/default`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ value: parseInt(value) }),
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const respData = await response.json();
      if (respData.status === "error") {
        addNotificationError(respData.data);
      }
    } catch (error) {
      console.error("Error posting data:", error);
    } finally {
      refresh();
    }
  };

  const handleOk = () => {
    setOpen(false);
    setDefault();
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
  refresh: PropTypes.func.isRequired,
};
