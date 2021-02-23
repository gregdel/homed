import React, { useState } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";
import { useParams } from "react-router-dom";

import { Form, Input, Button } from "antd";

import { setScheduleDefault } from "../../actions/schedule";

export const DefaultValue = ({ defaultValue: v }) => {
  const dispatch = useDispatch();
  const { componentId } = useParams();
  const [defaultValue, setDefaultValue] = useState(v);

  const update = () => {
    dispatch(setScheduleDefault(componentId, defaultValue));
  };

  return (
    <Form layout="inline" onFinish={update} onFinishFailed={update}>
      <Form.Item
        label="Default value"
        name="default_value"
        value={defaultValue}
        initialValue={defaultValue}
        onChange={(e) => {
          setDefaultValue(e.target.value);
        }}
      >
        <Input type="number" />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit">
          Update
        </Button>
      </Form.Item>
    </Form>
  );
};
DefaultValue.propTypes = {
  defaultValue: PropTypes.number.isRequired,
};
