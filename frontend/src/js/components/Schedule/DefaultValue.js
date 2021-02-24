import React, { useState } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";
import { useParams } from "react-router-dom";

import { Input } from "antd";

import { setScheduleDefault } from "../../actions/schedule";

export const DefaultValue = ({ defaultValue: v }) => {
  const dispatch = useDispatch();
  const { componentId } = useParams();
  const [defaultValue, setDefaultValue] = useState(v);

  const update = () => {
    dispatch(setScheduleDefault(componentId, defaultValue));
  };

  return (
    <div>
      <div>Default value</div>
      <div>
        <Input
          onChange={(e) => {
            setDefaultValue(e.target.value);
          }}
          onPressEnter={update}
          value={defaultValue}
          suffix="°C"
          type="number"
        />
      </div>
    </div>
  );
};
DefaultValue.propTypes = {
  defaultValue: PropTypes.number.isRequired,
};
